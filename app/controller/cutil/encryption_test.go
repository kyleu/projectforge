package cutil_test

import (
	"testing"

	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/lib/log"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	logger, _ := log.InitLogging(true)

	gcmTests := []struct{ plaintext string }{
		{plaintext: "Hello, world!"},
	}

	for _, tt := range gcmTests {
		ciphertext, err := cutil.EncryptMessage(tt.plaintext, logger)
		if err != nil {
			t.Fatal(err)
		}

		plaintext, err := cutil.DecryptMessage(ciphertext, logger)
		if err != nil {
			t.Fatal(err)
		}

		if plaintext != tt.plaintext {
			t.Errorf("plaintexts don't match")
		}
	}
}
