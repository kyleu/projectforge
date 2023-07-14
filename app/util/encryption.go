// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash/fnv"
	"io"
	"strings"

	"github.com/pkg/errors"
)

var _encryptKey string

func EncryptMessage(key []byte, message string, logger Logger) (string, error) {
	block, err := newCipher(key, logger)
	if err != nil {
		return "", errors.Wrap(err, "could not create new cipher")
	}

	byteMsg := []byte(message)
	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "could not encrypt")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(key []byte, message string, logger Logger) (string, error) {
	block, err := newCipher(key, logger)
	if err != nil {
		return "", errors.Wrap(err, "could not create new cipher")
	}

	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", errors.Wrap(err, "could not base64 decode")
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

func newCipher(key []byte, logger Logger) (cipher.Block, error) {
	if key == nil {
		key = getKey(logger)
	}
	for i := len(key); i < 16; i++ {
		key = append(key, ' ')
	}
	key = key[:16]
	return aes.NewCipher(key)
}

func HashFNV32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// HashSHA256 returns a Base64-encoded string representing the SHA-256 hash of the argument.
func HashSHA256(s string) string {
	h := sha256.New()
	ret := h.Sum([]byte(s))
	return base64.URLEncoding.EncodeToString(ret)
}

func getKey(logger Logger) []byte {
	if _encryptKey == "" {
		env := strings.ReplaceAll(AppKey, "-", "_") + "_encryption_key"
		_encryptKey = GetEnv(env)
		if _encryptKey == "" {
			if logger != nil {
				logger.Warnf("using default encryption key; set environment variable [%s] to save sessions between restarts", env)
			}
			_encryptKey = AppKey + "_secret"
		}
		for i := len(_encryptKey); i < 16; i++ {
			_encryptKey += " "
		}
		_encryptKey = _encryptKey[:16]
	}
	return []byte(_encryptKey)
}
