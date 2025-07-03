package client

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/lenarlenar/mygokeeper/internal/client/storage"
)

type Client struct {
	ServerURL string
	Token     string
	Reader    *bufio.Reader
}

func (c *Client) Register() {
	login := c.prompt("Login: ")
	password := c.prompt("Password: ")
	payload := map[string]string{
		"login":    login,
		"password": password,
	}
	_, err := doJSONRequest("POST", c.ServerURL+"/register", "", payload)
	if err != nil {
		fmt.Println("Request failed:", err)
	}
}

func (c *Client) Login() {
	login := c.prompt("Login: ")
	password := c.prompt("Password: ")
	payload := map[string]string{
		"login":    login,
		"password": password,
	}
	resp, err := doJSONRequest("POST", c.ServerURL+"/login", "", payload)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Login failed:", resp.Status)
		return
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Failed to decode:", err)
		return
	}
	if err := saveToken(result.Token); err != nil {
		fmt.Println("Failed to save token:", err)
		return
	}
	fmt.Println("Login successful.")
	c.Token = result.Token
}

func (c *Client) Add() {
	if c.Token == "" {
		fmt.Println("Please login first.")
		return
	}
	t := c.prompt("Type: ")
	d := c.prompt("Data: ")
	m := c.prompt("Meta: ")
	payload := map[string]string{
		"type": t,
		"data": base64.StdEncoding.EncodeToString([]byte(d)),
		"meta": m,
	}
	resp, err := doJSONRequest("POST", c.ServerURL+"/record", c.Token, payload)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		fmt.Println("Record saved.")
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

func (c *Client) Get() {
	if c.Token == "" {
		fmt.Println("Please login first.")
		return
	}
	resp, err := doJSONRequest("GET", c.ServerURL+"/records", c.Token, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		return
	}
	var records []struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		Data string `json:"data"`
		Meta string `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		fmt.Println("Failed to decode:", err)
		return
	}
	for _, r := range records {
		fmt.Printf("[%d] %s — %s\n", r.ID, r.Type, r.Meta)
		decoded, err := base64.StdEncoding.DecodeString(r.Data)
		if err != nil {
			fmt.Println("Failed to decode record:", err)
			continue
		}
		fmt.Println("   →", string(decoded))
	}
}

func (c *Client) Sync() {
	if c.Token == "" {
		fmt.Println("Please login first.")
		return
	}

	resp, err := doJSONRequest("GET", c.ServerURL+"/records", c.Token, nil)
	if err != nil {
		fmt.Println("Failed to get records from server")
		return
	}
	defer resp.Body.Close()
	var serverRecords []storage.LocalRecord
	if err := json.NewDecoder(resp.Body).Decode(&serverRecords); err != nil {
		fmt.Println("Decode failed:", err)
		return
	}

	localRecords, err := storage.LoadLocal()
	if err != nil {
		fmt.Println("Failed to load local records:", err)
		return
	}

	for _, r := range localRecords {
		resp, err := doJSONRequest("POST", c.ServerURL+"/record", c.Token, r)
		if err != nil {
			fmt.Println("Failed to upload record:", err)
			continue
		}
		resp.Body.Close()
	}

	combined := append(localRecords, serverRecords...)
	storage.SaveLocal(combined)
	fmt.Println("Sync completed.")
}

func (c *Client) Delete() {
	if c.Token == "" {
		fmt.Println("Please login first.")
		return
	}
	id := c.prompt("Enter record ID to delete: ")
	resp, err := doJSONRequest("DELETE", c.ServerURL+"/record/"+id, c.Token, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("Record deleted.")
	} else {
		fmt.Println("Delete failed:", resp.Status)
	}
}

func (c *Client) Update() {
	if c.Token == "" {
		fmt.Println("Please login first.")
		return
	}
	id := c.prompt("Enter record ID to update: ")
	typeVal := c.prompt("New Type: ")
	data := c.prompt("New Data: ")
	meta := c.prompt("New Meta: ")
	payload := map[string]string{
		"type": typeVal,
		"data": base64.StdEncoding.EncodeToString([]byte(data)),
		"meta": meta,
	}
	resp, err := doJSONRequest("PUT", c.ServerURL+"/record/"+id, c.Token, payload)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("Record updated.")
	} else {
		fmt.Println("Update failed:", resp.Status)
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

func LoadToken() (string, error) {
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

// doJSONRequest отправляет HTTP-запрос с JSON и токеном.
// method = POST, PUT, GET, DELETE
func doJSONRequest(method, url string, token string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}

func (c *Client) prompt(label string) string {
	fmt.Print(label)
	input, err := c.Reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}
