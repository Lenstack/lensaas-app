package utils

import (
	"fmt"
	"testing"
)

func TestNewCode(t *testing.T) {
	code := NewCode()
	fmt.Println(code)
	if code < 1000000 || code > 9999999 {
		t.Errorf("NewCode() = %v; want 1000000 <= code <= 9999999", code)
	}
}
