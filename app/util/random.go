// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/url"
	"time"

	"github.com/samber/lo"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	lo.ForEach(b, func(_ byte, i int) {
		b[i] = charset[RandomInt(len(charset))]
	})
	return string(b)
}

func RandomInt(maxExclusive int) int {
	maxCount := big.NewInt(int64(maxExclusive))
	ret, err := rand.Int(rand.Reader, maxCount)
	if err != nil {
		panic(errMsg(err))
	}
	return int(ret.Int64())
}

func RandomFloat(maxExclusive int) float64 {
	ret := RandomInt(maxExclusive * 1000)
	return float64(ret) / 1000
}

func RandomBool() bool {
	return RandomInt(2) == 0
}

func RandomValueMap(keys int) ValueMap {
	ret := ValueMap{}
	lo.Times(keys, func(_ int) struct{} {
		ret[RandomString(4)] = RandomString(10)
		return EmptyStruct
	})
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
	from := TimeCurrentUnix()
	to := TimeCurrent().AddDate(2, 0, 0).Unix()
	rnd := RandomInt(int(to - from))
	return time.Unix(from+int64(rnd), 0)
}

func RandomURL() *url.URL {
	u, _ := url.Parse("https://" + RandomString(6) + ".com/" + RandomString(6) + "/" + RandomString(8) + "." + RandomString(3))
	return u
}

func RandomElements[T any](l []T, x int) []T {
	if len(l) <= x {
		return lo.Shuffle(l)
	}
	return lo.Map(lo.Range(x), func(_ int, _ int) T {
		return l[RandomInt(len(l))]
	})
}

func RandomDiffs(size int) Diffs {
	return lo.Times(size, func(_ int) *Diff {
		return &Diff{Path: RandomString(8), Old: RandomString(12), New: RandomString(12)}
	})
}

func errMsg(err error) string {
	return fmt.Sprintf("source of randomness unavailable: %+v", err.Error())
}
