package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var currentConfig Config

func SetConfig(c Config) {
	currentConfig = c
}

func sendJSON(endpoint, login, password string) {
	payload := map[string]string{
		"login":    strings.TrimSpace(login),
		"password": strings.TrimSpace(password),
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(currentConfig.ServerURL+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Success:", resp.Status)
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

func saveToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := home + "/.gokeeper_token"
	return os.WriteFile(path, []byte(token), 0600)
}

func loadToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := home + "/.gokeeper_token"
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
