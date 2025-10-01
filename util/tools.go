package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"github.com/google/uuid"
)

var ToolsUtil = &toolsUtil{}

type toolsUtil struct{}

// md5
func (t *toolsUtil) Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// verifyMd5
func (t *toolsUtil) VerifyMd5(str, hash string) bool {
	return t.Md5(str) == hash
}

// uuid
func (t *toolsUtil) Uuid() string {
	return uuid.New().String()
}

// RandomString
func (t *toolsUtil) RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomNumber
func (t *toolsUtil) RandomNumber(n int) string {
	numbers := []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}
