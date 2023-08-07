// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
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
	max := big.NewInt(int64(maxExclusive))
	ret, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errMsg(err))
	}
	return int(ret.Int64())
}

func RandomFloat(maxExclusive int) float64 {
	max := big.NewInt(int64(maxExclusive) * 1000)
	ret, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errMsg(err))
	}
	return float64(ret.Int64()) / 1000
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
	lo.Times(keys, func(_ int) struct{} {
		ret[RandomString(4)] = RandomString(10)
		return struct{}{}
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

	rnd, err := rand.Int(rand.Reader, big.NewInt(to-from))
	if err != nil {
		panic(errMsg(err))
	}

	return time.Unix(from+rnd.Int64(), 0)
}

func RandomDiffs(size int) Diffs {
	return lo.Times(size, func(_ int) *Diff {
		return &Diff{Path: RandomString(8), Old: RandomString(12), New: RandomString(12)}
	})
}

func errMsg(err error) string {
	return fmt.Sprintf("source of randomness unavailable: %+v", err.Error())
}
