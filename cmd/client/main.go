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
	client.SetConfig(cfg)

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("GoKeeper CLI\nVersion: %s\nBuild date: %s\n", Version, BuildDate)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to GoKeeper CLI")
	fmt.Println("Available commands: register, login, add, get, delete, update, sync, exit")

	for {
		fmt.Print("> ")
		rawInput, _ := reader.ReadString('\n')
		input := strings.Fields(strings.TrimSpace(rawInput))
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "register":
			client.HandleRegister(reader)
		case "login":
			client.HandleLogin(reader)
		case "add":
			client.HandleAdd(reader)
		case "get":
			client.HandleGet()
		case "sync":
			client.HandleSync()
		case "delete":
			client.HandleDelete(reader)
		case "update":
			client.HandleUpdate(reader)
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command:", input[0])
		}
	}
}
