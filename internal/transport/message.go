// internal/transport/message.go
package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func SendMessage(conn net.Conn, message string) error {
	_, err := fmt.Fprintln(conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func ReadMessage(reader *bufio.Reader) (string, error) {
	message, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", ErrConnectionClosed
		}
		return "", fmt.Errorf("error reading message: %w", err)
	}
	return strings.TrimSpace(message), nil
}

type MessageProtocol struct {
	Delimiter byte
}

func NewMessageProtocol(delimiter byte) *MessageProtocol {
	return &MessageProtocol{Delimiter: delimiter}
}

func (mp *MessageProtocol) Encode(message string) []byte {
	return append([]byte(message), mp.Delimiter)
}

func (mp *MessageProtocol) Decode(data []byte) string {
	return strings.TrimSuffix(string(data), string(mp.Delimiter))
}

var (
	ErrConnectionClosed = fmt.Errorf("connection closed by peer")
)
