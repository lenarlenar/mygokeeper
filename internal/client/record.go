package client

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleAdd(r *bufio.Reader) {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Please login first.")
		return
	}

	fmt.Print("Type (password/card/text/binary): ")
	t, _ := r.ReadString('\n')
	fmt.Print("Data: ")
	d, _ := r.ReadString('\n')
	fmt.Print("Meta: ")
	m, _ := r.ReadString('\n')

	payload := map[string]string{
		"type": strings.TrimSpace(t),
		"data": base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(d))),
		"meta": strings.TrimSpace(m),
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", currentConfig.ServerURL+"/record", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Record saved.")
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

func HandleGet() {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Please login first.")
		return
	}

	req, _ := http.NewRequest("GET", currentConfig.ServerURL+"/records", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	var records []struct {
		Type string `json:"type"`
		Data string `json:"data"`
		Meta string `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		fmt.Println("Failed to decode response:", err)
		return
	}

	for i, r := range records {
		fmt.Printf("[%d] %s — %s\n", i+1, r.Type, r.Meta)
		data, _ := base64.StdEncoding.DecodeString(r.Data)
		fmt.Println("   →", string(data))
	}
}

func HandleDelete(reader *bufio.Reader) {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Please login first.")
		return
	}

	fmt.Print("Enter record ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	req, _ := http.NewRequest("DELETE", currentConfig.ServerURL+"/record/"+idStr, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Record deleted.")
	} else {
		fmt.Println("Delete failed:", resp.Status)
	}
}

func HandleUpdate(reader *bufio.Reader) {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Please login first.")
		return
	}

	fmt.Print("Enter record ID to update: ")
	idStr, _ := reader.ReadString('\n')

	fmt.Print("New Type: ")
	t, _ := reader.ReadString('\n')
	fmt.Print("New Data: ")
	d, _ := reader.ReadString('\n')
	fmt.Print("New Meta: ")
	m, _ := reader.ReadString('\n')

	payload := map[string]string{
		"type": strings.TrimSpace(t),
		"data": base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(d))),
		"meta": strings.TrimSpace(m),
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", currentConfig.ServerURL+"/record/"+strings.TrimSpace(idStr), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Record updated.")
	} else {
		fmt.Println("Update failed:", resp.Status)
	}
}
