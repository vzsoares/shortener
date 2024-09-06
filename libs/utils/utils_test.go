package utils

import (
	"testing"
)

func TestUtils(t *testing.T) {
	result := Utils("works")
	if result != "Utils works" {
		t.Error("Expected Utils to append 'works'")
	}
}
