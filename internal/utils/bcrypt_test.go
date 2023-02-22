package utils

import "testing"

func TestComparePassword(t *testing.T) {
	b := &Bcrypt{}
	password := "password123"
	hashedPassword, err := b.HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	// Test valid password comparison
	err = b.ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Valid password comparison failed: %v", err)
	}

	// Test invalid password comparison
	err = b.ComparePassword(hashedPassword, "wrong-password")
	if err == nil {
		t.Errorf("Invalid password comparison should return an error")
	}

}

func TestHashPassword(t *testing.T) {
	b := &Bcrypt{}
	password := "password123"
	hashedPassword, err := b.HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if hashedPassword == password {
		t.Errorf("Hashed password should not be equal to the original password")
	}
}
