package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/SaeedMPro/socket-chat/model"
	"github.com/SaeedMPro/socket-chat/util"
)

var (
	clientConfig model.Config
	currentConn  net.Conn
	connMux      sync.Mutex
)

func init() {
	if err := util.LoadConfig("./config.json", &clientConfig); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	util.ClearScreen()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [client-one|client-two]")
		return
	}
	
	clientType := os.Args[1]
	self, other := getClientConfig(clientType)

	go startListener(self)
	connectToPeer(other)
	
	select {} 
}

func getClientConfig(clientType string) (model.Client, model.Client) {
	switch clientType {
	case "client-one":
		return clientConfig.ClientOne, clientConfig.ClientTwo
	case "client-two":
		return clientConfig.ClientTwo, clientConfig.ClientOne
	default:
		fmt.Println("Invalid client type")
		os.Exit(1)
		return model.Client{}, model.Client{}
	}
}

func startListener(self model.Client) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", self.Host, self.Port))
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			continue
		}
		handleNewConnection(connection)
	}

}

func handleNewConnection(conn net.Conn) {
	connMux.Lock()
	defer connMux.Unlock()

	if currentConn != nil {
		currentConn.Close()
	}
	currentConn = conn

	//TODO : implement send and receive message
}

func connectToPeer(peer model.Client) {
	peerAddr := fmt.Sprintf("%v:%v", peer.Host, peer.Port)
	
	for {
		connection, err := net.Dial("tcp", peerAddr)
		if err != nil {
			fmt.Printf("Connection attempt failed: %v\n", err)
			continue
		}
		handleNewConnection(connection)
		break
	}
}