package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// NIKCrypto provides encryption, decryption, and hashing operations for NIK (ID numbers).
// Uses AES-256-GCM for encryption and SHA-256 for hashing.
type NIKCrypto struct {
	key []byte
}

// NewNIKCrypto creates a new NIKCrypto instance from a hex-encoded 32-byte key.
func NewNIKCrypto(hexKey string) (*NIKCrypto, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("crypto.NewNIKCrypto: invalid hex key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("crypto.NewNIKCrypto: key must be 32 bytes, got %d", len(key))
	}
	return &NIKCrypto{key: key}, nil
}

// Encrypt encrypts a NIK using AES-256-GCM.
// Returns the encrypted data as a byte slice (nonce prepended to ciphertext).
func (n *NIKCrypto) Encrypt(plainNIK string) ([]byte, error) {
	block, err := aes.NewCipher(n.key)
	if err != nil {
		return nil, fmt.Errorf("crypto.Encrypt: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto.Encrypt: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("crypto.Encrypt: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainNIK), nil)
	return ciphertext, nil
}

// Decrypt decrypts a NIK that was encrypted with Encrypt.
func (n *NIKCrypto) Decrypt(encrypted []byte) (string, error) {
	block, err := aes.NewCipher(n.key)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(encrypted) < nonceSize {
		return "", fmt.Errorf("crypto.Decrypt: ciphertext too short")
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: %w", err)
	}

	return string(plaintext), nil
}

// Hash returns the SHA-256 hex digest of a NIK for index-based lookups.
func Hash(plainNIK string) string {
	h := sha256.Sum256([]byte(plainNIK))
	return hex.EncodeToString(h[:])
}
