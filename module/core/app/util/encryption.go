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
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	_encryptKey     string
	_encryptKeyOnce sync.Once
)

func EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncryptMessage(key []byte, message string, logger Logger) (string, error) {
	cipherText, err := EncryptBytes(key, []byte(message), logger)
	if err != nil {
		return "", err
	}
	return EncodeBytes(cipherText), nil
}

func EncryptBytes(key []byte, b []byte, logger Logger) ([]byte, error) {
	block, err := newCipher(key, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new cipher")
	}

	if len(b) > 1024*1024*64 {
		return nil, errors.New("message is too large")
	}
	cipherText := make([]byte, aes.BlockSize+len(b))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(err, "could not encrypt")
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], b)

	return cipherText, nil
}

func DecryptBytes(key []byte, b []byte, logger Logger) ([]byte, error) {
	block, err := newCipher(key, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new cipher")
	}

	if len(b) < aes.BlockSize {
		return nil, errors.New("invalid ciphertext block size")
	}

	iv := b[:aes.BlockSize]
	cipherText := b[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func DecryptMessage(key []byte, message string, logger Logger) (string, error) {
	cipherText, err := DecodeBytes(message)
	if err != nil {
		return "", errors.Wrap(err, "could not base64 decode")
	}
	byteMsg, err := DecryptBytes(key, cipherText, logger)
	if err != nil {
		return "", err
	}
	return string(byteMsg), nil
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

func HashFNV128UUID(s string) uuid.UUID {
	h := fnv.New128a()
	_, _ = h.Write([]byte(s))
	b := h.Sum(make([]byte, 0, 16))
	ret, _ := uuid.FromBytes(b)
	return ret
}

// HashSHA256 returns a Base64-encoded string representing the SHA-256 hash of the argument.
func HashSHA256(s string) string {
	h := sha256.New()
	ret := h.Sum([]byte(s))
	return base64.URLEncoding.EncodeToString(ret)
}

func getKey(logger Logger) []byte {
	_encryptKeyOnce.Do(func() {
		envName := strings.ReplaceAll(AppKey, "-", "_") + "_encryption_key"
		_encryptKey = GetEnv(envName)
		if _encryptKey == "" {
			if logger != nil {
				logger.Warnf("using default encryption key; set environment variable [%s] to save sessions between restarts", envName)
			}
			_encryptKey = AppKey + "_secret"
		}
		for i := len(_encryptKey); i < 16; i++ {
			_encryptKey += " "
		}
		_encryptKey = _encryptKey[:16]
	})
	return []byte(_encryptKey)
}
