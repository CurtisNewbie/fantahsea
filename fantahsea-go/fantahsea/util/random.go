package util

import (
	"math/rand"
)

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generate a random sequence number with specified prefix
func GenNo(prefix string) string {
	return prefix + randStr(10)
}
