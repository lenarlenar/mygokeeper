package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lenarlenar/mygokeeper/internal/client"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
)

func main() {
	cfg := client.Load()

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("GoKeeper CLI\nVersion: %s\nBuild date: %s\n", Version, BuildDate)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	token, _ := client.LoadToken()

	c := &client.Client{
		ServerURL: cfg.ServerURL,
		Token:     token,
		Reader:    reader,
	}
	fmt.Println("Welcome to GoKeeper CLI")
	fmt.Println("Available commands: register, login, add, get, delete, update, sync, exit")

	for {
		fmt.Print("> ")
		rawInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		input := strings.Fields(strings.TrimSpace(rawInput))
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "register":
			c.Register()
		case "login":
			c.Login()
		case "add":
			c.Add()
		case "get":
			c.Get()
		case "sync":
			c.Sync()
		case "delete":
			c.Delete()
		case "update":
			c.Update()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command:", input[0])
		}
	}
}
