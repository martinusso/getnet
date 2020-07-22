package getnet

import (
	"testing"
)

func TestSubstring(t *testing.T) {
	value := "Hello, playground"

	got := maxLength(value, 50)
	if got != value {
		t.Errorf("Expected '%s', got '%s'", value, got)
	}

	got = maxLength(value, 6)
	expected := "Hello,"
	if got != expected {
		t.Errorf("Expected '%s', got '%s'", expected, got)
	}
}
