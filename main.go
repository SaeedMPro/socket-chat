package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/SaeedMPro/socket-chat/config"
	"github.com/SaeedMPro/socket-chat/model"
	"github.com/SaeedMPro/socket-chat/util"
)

var (
	clientConfig model.Config
	currentConn  net.Conn
)

func init() {
	if err := config.LoadConfig("./config/config.json", &clientConfig); err != nil {
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
	self, other := config.GetClientConfig(clientType, clientConfig)

	go startListener(self)
	go connectToPeer(other)

	select {}
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
			time.Sleep(1 * time.Second)
			continue
		}
		handleNewConnection(connection)
	}

}

func connectToPeer(peer model.Client) {
	peerAddr := peer.Address()

	for {
		connection, err := net.Dial("tcp", peerAddr)
		if err != nil {
			fmt.Printf("Connection attempt failed: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		handleNewConnection(connection)
		break
	}
}

func handleNewConnection(conn net.Conn) {
	if currentConn != nil {
		currentConn.Close()
	}
	currentConn = conn
	util.ClearScreen()

	fmt.Println("-----------------------------Connection created-----------------------------")
	fmt.Println("  <YOU>\t\t\t\t\t\t\t\t<Peer>")

	done := make(chan struct{})
	go sendMessages(conn, done)
	go receiveMessages(conn, done)
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
