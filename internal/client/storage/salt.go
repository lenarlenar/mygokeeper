package storage

import (
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
)

func loadOrCreateSalt() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".gokeeper_salt")

	if data, err := os.ReadFile(path); err == nil {
		return data, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	// Генерируем новую соль
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	err = os.WriteFile(path, salt, 0600)
	return salt, err
}
