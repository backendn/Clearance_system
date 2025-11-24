package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteRune(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func RandomEmail() string {
	return RandomString(8) + "@example.com"
}

func RandomStudentNumber() string {
	// Example simple format — change if needed
	return fmt.Sprintf("%s/%d", RandomString(3), rand.Intn(9999))
}

func RandomPhone() string {
	// Ethiopian-style numbers: 09xxxxxxxx
	return fmt.Sprintf("09%08d", rand.Intn(100000000))
}
