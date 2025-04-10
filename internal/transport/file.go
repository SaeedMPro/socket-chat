package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleFileTransfer(header string, reader *bufio.Reader) {
	parts := strings.SplitN(header, " ", 3)
	if len(parts) != 3 {
		fmt.Println("\nInvalid file header")
		return
	}

	filename := parts[1]
	size, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("\nInvalid file size:", parts[2])
		return
	}

	data := make([]byte, size)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		fmt.Printf("\nError receiving file: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("\nError saving file: %v\n", err)
	} else {
		fmt.Printf("\nFile received: %s (%d bytes)\n", filename, size)
	}
}

func SendFile(conn net.Conn, filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
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
}
