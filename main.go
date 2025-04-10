package main

import (
	"bufio"
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

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [client-one|client-two]")
		os.Exit(1)
	}

	util.ClearScreen()
}

func main() {
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
		panic("Invalid argument input!!")
	}
}

func startListener(self model.Client) {
	listener, err := net.Listen("tcp", self.Address())
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
	util.ClearScreen()

	fmt.Println("-----------------------------Connection created-----------------------------")
	fmt.Println("  <YOU>\t\t\t\t\t\t\t\t<Peer>")

	done := make(chan struct{})
	go receiveMessages(conn, done)
	go sendMessages(conn, done)
}

func connectToPeer(peer model.Client) {
	peerAddr := peer.Address()

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

func receiveMessages(conn net.Conn, done chan struct{}) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nConnection lost")
			break
		}
		fmt.Printf("\n\t\t\t\t\t\t\t\t %s ", msg)
	}

	close(done)
}

func sendMessages(conn net.Conn, done chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-done:
			return
		default:
			fmt.Print("\n ")
			if !scanner.Scan() {
				return
			}
			msg := scanner.Text()

			_, err := fmt.Fprintln(conn, msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}
