package client

import (
	"os"
)

type Config struct {
	ServerURL string
}

func Load() Config {
	serverURL := os.Getenv("GOKEEPER_SERVER")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	return Config{
		ServerURL: serverURL,
	}
}
