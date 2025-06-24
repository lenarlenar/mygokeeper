package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lenarlenar/mygokeeper/internal/client/storage"
)

func HandleSync() {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Please login first.")
		return
	}

	// 1. Получаем записи с сервера
	req, _ := http.NewRequest("GET", currentConfig.ServerURL+"/records", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get records from server")
		return
	}
	defer resp.Body.Close()

	var serverRecords []storage.LocalRecord
	if err := json.NewDecoder(resp.Body).Decode(&serverRecords); err != nil {
		fmt.Println("Decode failed:", err)
		return
	}

	// 2. Получаем локальные записи
	localRecords, _ := storage.LoadLocal()

	// 3. Отправляем локальные записи на сервер
	for _, r := range localRecords {
		body, _ := json.Marshal(r)
		req, _ := http.NewRequest("POST", currentConfig.ServerURL+"/record", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}

	// 4. Объединяем и сохраняем в файл
	combined := append(localRecords, serverRecords...)
	storage.SaveLocal(combined)

	fmt.Println("Sync completed.")
}
