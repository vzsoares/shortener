package utils

import (
	"errors"
	"fmt"
	"math/rand"
)

var UriSafeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~()'!*:@,;"

func GenUriSafeRash(size int, prefix string) string {
	l := len(UriSafeChars)
	s := ""
	for range size {
		rd := rand.Intn(l)
		c := UriSafeChars[rd]
		s += string(c)
	}
	return fmt.Sprintf("%v%v", prefix, s)
}

func RetryN[T any](fn func(i int) (T, error), count int) (T, error) {
	for i := range count {
		v, err := fn(i + 1)
		if err == nil {
			return v, nil
		}
	}
	var v T
	return v, errors.New("Failed miserably")
}
