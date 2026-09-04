package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "StrongPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Fatal("password must not be stored as plain text")
	}

	if !VerifyPassword(password, hash) {
		t.Fatal("VerifyPassword() failed for correct password")
	}

	if VerifyPassword("WrongPassword123!", hash) {
		t.Fatal("VerifyPassword() accepted incorrect password")
	}
}

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	password := "StrongPassword123!"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword() error = %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("same password should produce different hashes because of random salts")
	}
}
