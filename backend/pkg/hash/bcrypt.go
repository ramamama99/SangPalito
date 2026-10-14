package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost is the cost factor used for bcrypt hashing. Minimum 12, jangan kurang.
const BcryptCost = 12

// HashPassword menghasilkan hash bcrypt yang sudah include salt otomatis.
// bcrypt secara internal men-generate random salt — JANGAN tambah salt manual.
func HashPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash.HashPassword: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword membandingkan plain password dengan hash.
// Mengembalikan nil jika cocok, error jika tidak.
func VerifyPassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
