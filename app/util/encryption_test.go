//go:build test_all || !func_test

// Content managed by Project Forge, see [projectforge.md] for details.
package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

var hashTests = []struct {
	plaintext string
	hash32    uint32
	hashSHA   string
}{
	{plaintext: "Hello, world!", hash32: 3985698964, hashSHA: "SGVsbG8sIHdvcmxkIeOwxEKY_BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhV"},
	{plaintext: "Goodbye, cruel world!", hash32: 55824456, hashSHA: "R29vZGJ5ZSwgY3J1ZWwgd29ybGQh47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="},
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	for _, tt := range hashTests {
		ciphertext, err := util.EncryptMessage(nil, tt.plaintext, nil)
		if err != nil {
			t.Fatal(err)
		}

		plaintext, err := util.DecryptMessage(nil, ciphertext, nil)
		if err != nil {
			t.Fatal(err)
		}

		if plaintext != tt.plaintext {
			t.Errorf("plaintexts don't match")
		}
	}
}

func TestHashFNV32(t *testing.T) {
	t.Parallel()

	for _, tt := range hashTests {
		h := util.HashFNV32(tt.plaintext)
		if h != tt.hash32 {
			t.Errorf("FNC32 hash didn't match for [%s], expected [%d], observed [%d]", tt.plaintext, tt.hash32, h)
		}
	}
}

func TestHashSHA256(t *testing.T) {
	t.Parallel()

	for _, tt := range hashTests {
		h := util.HashSHA256(tt.plaintext)
		if h != tt.hashSHA {
			t.Errorf("SHA-256 hash didn't match for [%s], expected [%s], observed [%s]", tt.plaintext, tt.hashSHA, h)
		}
	}
}
