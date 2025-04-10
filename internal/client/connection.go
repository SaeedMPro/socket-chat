package client

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/SaeedMPro/socket-chat/internal/transport"
	"github.com/SaeedMPro/socket-chat/internal/ui"
)

func (c *ChatClient) receiveMessages() {
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				ui.ShowMessage("\nPeer disconnected")
			}
			break
		}
		msg = strings.TrimSpace(msg)

		if strings.HasPrefix(msg, "FILE ") {
			transport.HandleFileTransfer(msg, reader)
		} else {
			ui.DisplayPeerMessage(msg)
		}
	}
}

func (c *ChatClient) sendMessages() {
	defer c.conn.Close()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		input := ui.GetUserInput(scanner)

		if strings.HasPrefix(input, "/sendfile ") {
			parts := strings.SplitN(input, " ", 2)
			if len(parts) < 2 {
				ui.ShowMessage("Usage: /sendfile <filename>")
				continue
			}
			transport.SendFile(c.conn, parts[1])
		} else {
			transport.SendMessage(c.conn, input)
		}
	}
}
