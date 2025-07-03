package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/term"
)

type LocalRecord struct {
	Type string `json:"type"`
	Data string `json:"data"` // base64 string
	Meta string `json:"meta"`
}

func getStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gokeeper_data.json"), nil
}

func LoadLocal() ([]LocalRecord, error) {
	path, err := getStoragePath()
	if err != nil {
		return nil, err
	}

	pass := promptPassword("Enter passphrase to unlock local storage: ")
	return LoadRecordsFrom(path, pass)
}

func SaveLocal(records []LocalRecord) error {
	path, err := getStoragePath()
	if err != nil {
		return err
	}

	pass := promptPassword("Enter passphrase for encryption: ")
	return SaveRecordsTo(path, pass, records)
}

func promptPassword(prompt string) string {
	fmt.Print(prompt)
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	return string(pass)
}

// SaveRecordsTo шифрует и сохраняет записи в указанный файл.
func SaveRecordsTo(path string, passphrase string, records []LocalRecord) error {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	enc, err := encrypt(jsonData, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(path, enc, 0600)
}

// LoadRecordsFrom расшифровывает и загружает записи из файла.
func LoadRecordsFrom(path string, passphrase string) ([]LocalRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dec, err := decrypt(data, passphrase)
	if err != nil {
		return nil, err
	}

	var records []LocalRecord
	err = json.Unmarshal(dec, &records)
	return records, err
}
