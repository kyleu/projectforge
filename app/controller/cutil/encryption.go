package cutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/util"
)

var key string

func EncryptMessage(message string, logger *zap.SugaredLogger) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(getKey(logger))
	if err != nil {
		return "", errors.Wrap(err, "could not create new cipher")
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "could not encrypt")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(message string, logger *zap.SugaredLogger) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", errors.Wrap(err, "could not base64 decode")
	}

	block, err := aes.NewCipher(getKey(logger))
	if err != nil {
		return "", errors.Wrap(err, "could not create new cipher")
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func getKey(logger *zap.SugaredLogger) []byte {
	if key == "" {
		env := strings.ReplaceAll(util.AppKey, "-", "_") + "_encryption_key"
		key = os.Getenv(env)
		if key == "" {
			logger.Warnf("using default encryption key\nset environment variable [%s] to save sessions between restarts", env)
			key = util.AppKey + "_secret"
		}
		for i := len(key); i < 16; i++ {
			key += " "
		}
		key = key[:16]
	}
	return []byte(key)
}
