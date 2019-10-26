package random

import (
	"math/rand"
	"time"
)

func RandString(size int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
