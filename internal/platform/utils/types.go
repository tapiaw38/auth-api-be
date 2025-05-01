package utils

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func ToDateString(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}

func ToInt(r rune) int {
	return int(r - '0')
}

func ToPointer[T any](value T) *T {
	return &value
}

func GetEncodedString() (string, error) {
	size := 16
	strBytes := make([]byte, size)

	_, err := cryptorand.Read(strBytes)
	if err != nil {
		return "", err
	}

	strEncoded := hex.EncodeToString(strBytes)

	return strEncoded, nil
}

func RandomString(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, lenght)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
