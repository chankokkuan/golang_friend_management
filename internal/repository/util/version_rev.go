package util

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateVersionRev ...
func GenerateVersionRev(seqNumber int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%d-%s", seqNumber, string(b))
}
