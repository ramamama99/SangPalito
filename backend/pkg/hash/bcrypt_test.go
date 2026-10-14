package hash

import (
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := "SecretPassword123!"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if len(hashed) == 0 {
		t.Fatal("expected non-empty hash string")
	}

	if hashed == password {
		t.Fatal("hash should not match plain password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "SecurePass@2026"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	// Correct password must verify successfully
	if err := VerifyPassword(password, hashed); err != nil {
		t.Errorf("expected password to verify successfully, got: %v", err)
	}

	// Wrong password must fail
	wrongPassword := "WrongPassword999!"
	if err := VerifyPassword(wrongPassword, hashed); err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}
}
