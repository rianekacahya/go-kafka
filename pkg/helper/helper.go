package helper

import (
	"context"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

// RandomString function for random string
func RandomString(length int) string {
	var b strings.Builder

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func GetContextString(ctx context.Context, key string) string {
	if t, ok := ctx.Value(key).(string); ok {
		return t
	}
	return ""
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
