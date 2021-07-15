package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		idx := RandomInt(len(charset))
		b[i] = charset[idx]
	}
	return string(b)
}

func RandomInt(maxInclusive int) int {
	max := big.NewInt(int64(maxInclusive))
	ret, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errMsg(err))
	}
	return int(ret.Int64())
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(errMsg(err))
	}
	return b
}

func errMsg(err error) string {
	return fmt.Sprintf("source of randomness unavailable: %+v", err.Error())
}
