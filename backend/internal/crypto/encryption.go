package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"sync"
)

var (
	encryptionKey []byte
	keyOnce       sync.Once
	keyErr        error

	ErrKeyNotConfigured = errors.New("ASC_ENCRYPTION_KEY not configured")
	ErrInvalidKeyLength = errors.New("encryption key must be 32 bytes (64 hex characters)")
	ErrDecryptionFailed = errors.New("decryption failed")
)

func initKey() {
	keyHex := os.Getenv("ASC_ENCRYPTION_KEY")
	if keyHex == "" {
		keyErr = ErrKeyNotConfigured
		return
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		keyErr = err
		return
	}

	if len(key) != 32 {
		keyErr = ErrInvalidKeyLength
		return
	}

	encryptionKey = key
}

func getKey() ([]byte, error) {
	keyOnce.Do(initKey)
	if keyErr != nil {
		return nil, keyErr
	}
	return encryptionKey, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func Encrypt(plaintext []byte) (ciphertext, nonce []byte, err error) {
	key, err := getKey()
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	key, err := getKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// GenerateKey generates a random 32-byte encryption key and returns it as hex
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
