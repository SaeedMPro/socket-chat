package main

import (
	"fmt"
	"os"

	"github.com/SaeedMPro/socket-chat/config"
	"github.com/SaeedMPro/socket-chat/internal/client"
	"github.com/SaeedMPro/socket-chat/model"
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [client-one|client-two]")
		os.Exit(1)
	}
}

func main() {
	var clientConfig model.Config

	clientType := os.Args[1]
	err := config.LoadConfig("./config/config.json", &clientConfig)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	client.New(clientType, clientConfig).Run()
}
