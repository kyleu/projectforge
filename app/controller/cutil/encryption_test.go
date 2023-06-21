package cutil_test

import (
	"github.com/samber/lo"
	"testing"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/log"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	logger, _ := log.InitLogging(true)

	gcmTests := []struct{ plaintext string }{
		{plaintext: "Hello, world!"},
	}

	lo.ForEach(gcmTests, func(tt struct{ plaintext string }, _ int) {
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
	})
}
