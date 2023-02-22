package utils

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	j := Jwt{}
	token, err := j.GenerateToken("123", 500000)
	if err != nil {
		return
	}

	t.Log(token)
}

func TestValidateToken(t *testing.T) {
	j := Jwt{}
	token, err := j.GenerateToken("123", 500000)
	if err != nil {
		return
	}

	userId, err := j.ValidateToken(token)
	if err != nil {
		return
	}

	if userId != "123" {
		t.Error("User id is not the same")
	}

	t.Log(userId)
}
