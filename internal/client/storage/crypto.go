package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(passphrase string) []byte {
	salt, err := loadOrCreateSalt()
	if err != nil {
		panic("cannot load salt: " + err.Error())
	}
	return pbkdf2.Key([]byte(passphrase), salt, 100_000, 32, sha256.New)
}

// encrypt шифрует переданные данные с использованием пароля (PBKDF2 + AES-GCM).
func encrypt(data []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt расшифровывает данные, зашифрованные encrypt с тем же паролем.
func decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
