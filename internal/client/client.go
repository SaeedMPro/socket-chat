package client

import (
	"fmt"
	"net"
	"time"

	"github.com/SaeedMPro/socket-chat/config"
	"github.com/SaeedMPro/socket-chat/internal/ui"
	"github.com/SaeedMPro/socket-chat/model"
)

type ChatClient struct {
	Self      model.Client
	Peer      model.Client
	conn      net.Conn
	connState chan struct{}
}

func New(clientType string, cfg model.Config) *ChatClient {
	self, peer := config.GetClientConfig(clientType, cfg)
	return &ChatClient{
		Self:      self,
		Peer:      peer,
		connState: make(chan struct{}, 1),
	}
}

func (c *ChatClient) Run() {
	go c.startListener()
	go c.connectToPeer()

	select {}
}

func (c *ChatClient) startListener() {
	listener, err := net.Listen("tcp", c.Self.Address())
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.handleNewConnection(conn)
	}
}

func (c *ChatClient) connectToPeer() {
	for {
		conn, err := net.Dial("tcp", c.Peer.Address())
		if err != nil {
			fmt.Printf("Connection attempt failed: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.handleNewConnection(conn)
		break
	}
}

func (c *ChatClient) handleNewConnection(conn net.Conn) {
	if c.conn != nil {
		c.conn.Close()
		<-c.connState
	}

	c.conn = conn
	c.connState <- struct{}{}

	ui.InitCliUI()

	go c.receiveMessages()
	go c.sendMessages()
}
