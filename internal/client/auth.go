package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleRegister(r *bufio.Reader) {
	fmt.Print("Login: ")
	login, _ := r.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := r.ReadString('\n')

	sendJSON("/register", login, password)
}

func HandleLogin(r *bufio.Reader) {
	fmt.Print("Login: ")
	login, _ := r.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := r.ReadString('\n')

	payload := map[string]string{
		"login":    strings.TrimSpace(login),
		"password": strings.TrimSpace(password),
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(currentConfig.ServerURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Login failed:", resp.Status)
		return
	}

	var respBody struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		fmt.Println("Failed to parse token:", err)
		return
	}

	err = saveToken(respBody.Token)
	if err != nil {
		fmt.Println("Failed to save token:", err)
		return
	}

	fmt.Println("Login successful. Token saved.")
}
