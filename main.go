package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SaeedMPro/socket-chat/config"
	"github.com/SaeedMPro/socket-chat/model"
	"github.com/SaeedMPro/socket-chat/util"
)

var (
	clientConfig model.Config
	currentConn  net.Conn
	connMux      sync.Mutex
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
		go handleNewConnection(connection)
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
		go handleNewConnection(connection)
		break
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
	go sendMessages(conn, done)
	go receiveMessages(conn, done)
}

func receiveMessages(conn net.Conn, done chan struct{}) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nConnection lost")
			break
		}
		msg = strings.TrimSpace(msg)

		if strings.HasPrefix(msg, "FILE ") {
			parts := strings.SplitN(msg, " ", 3)
			if len(parts) != 3 {
				fmt.Println("\nInvalid file header")
				continue
			}

			filename := parts[1]
			size, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("\nInvalid file size:", parts[2])
				continue
			}

			data := make([]byte, size)
			_, err = io.ReadFull(reader, data)
			if err != nil {
				fmt.Printf("\nError receiving file: %v\n", err)
				continue
			}

			err = os.WriteFile(filename, data, 0644)
			if err != nil {
				fmt.Printf("\nError saving file: %v\n", err)
			} else {
				fmt.Printf("\nFile received: %s (%d bytes)\n", filename, size)
			}
		} else {
			fmt.Printf("\n\t\t\t\t\t\t\t\t %s\n", msg)
		}
	}

	close(done)
}

func sendMessages(conn net.Conn, done chan struct{}) {
	defer conn.Close()
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
			input := scanner.Text()

			if strings.HasPrefix(input, "/sendfile ") {
				parts := strings.SplitN(input, " ", 2)
				if len(parts) < 2 {
					fmt.Println("Usage: /sendfile <filename>")
					continue
				}

				filename := parts[1]
				data, err := os.ReadFile(filename)
				if err != nil {
					fmt.Printf("Error reading file: %v\n", err)
					continue
				}

				header := fmt.Sprintf("FILE %s %d\n", filepath.Base(filename), len(data))
				_, err = fmt.Fprint(conn, header)
				if err != nil {
					fmt.Printf("Error sending file header: %v\n", err)
					return
				}

				_, err = conn.Write(data)
				if err != nil {
					fmt.Printf("Error sending file data: %v\n", err)
					return
				}

				fmt.Printf("File '%s' sent successfully!\n", filename)
			} else {
				_, err := fmt.Fprintln(conn, input)
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
			}
		}
	}
}