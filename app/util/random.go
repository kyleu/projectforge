package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"time"
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

func RandomInt(maxExclusive int) int {
	max := big.NewInt(int64(maxExclusive))
	ret, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errMsg(err))
	}
	return int(ret.Int64())
}

func RandomBool() bool {
	ret, err := rand.Int(rand.Reader, big.NewInt(1))
	if err != nil {
		panic(errMsg(err))
	}
	return ret.Int64() == 1
}

func RandomValueMap(keys int) ValueMap {
	ret := ValueMap{}
	for i := 0; i < keys; i++ {
		ret[RandomString(4)] = RandomString(10)
	}
	return ret
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(errMsg(err))
	}
	return b
}

func RandomDate() time.Time {
	from := time.Now().Unix()
	to := time.Now().AddDate(2, 0, 0).Unix()

	rnd, err := rand.Int(rand.Reader, big.NewInt(to-from))
	if err != nil {
		panic(errMsg(err))
	}

	return time.Unix(from+rnd.Int64(), 0)
}

func errMsg(err error) string {
	return fmt.Sprintf("source of randomness unavailable: %+v", err.Error())
}
