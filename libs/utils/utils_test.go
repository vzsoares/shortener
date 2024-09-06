package utils

import (
	"errors"
	"strings"
	"testing"
)

func TestGenUriSafeRash(t *testing.T) {
	result := GenUriSafeRash(5, "")
	if len(result) != 5 {
		t.Error("Expected rash to have len of 5")
	}
	result = GenUriSafeRash(5, "a-")
	if !strings.HasPrefix(result, "a-") {
		t.Error("Expected rash to have prefix a-")
	}
	for _, c := range result {
		ok := strings.Contains(UriSafeChars, string(c))
		if !ok {
			t.Error("Expected char to be uri safe")
		}
	}
}

func TestRetryN(t *testing.T) {
	fn := func(i int) (string, error) {
		if i == 5 {
			return "Yeah", nil
		}
		return "", errors.New("Test")
	}
	_, err := RetryN(fn, 5)
	if err != nil {
		t.Error("Expected to run 5 times")
	}
}
