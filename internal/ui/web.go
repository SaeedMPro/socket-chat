package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	messages     []string
	mu           sync.Mutex
	connSend     func(string)
	connSendFile func(string)
)

func AddMessage(msg string) {
	mu.Lock()
	defer mu.Unlock()
	messages = append(messages, msg)
}

func ServeWebUI(addr string, sendMsg func(string), sendFile func(string)) {
	connSend = sendMsg
	connSendFile = sendFile

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/ui/index.html")
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		for _, msg := range messages {
			io.WriteString(w, msg+"\n")
		}
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		msg := r.FormValue("msg")
		if msg != "" {
			AddMessage("YOU: " + msg)
			if connSend != nil {
				connSend(msg)
			}
		}
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("file")
		defer file.Close()

		tmpPath := "./" + header.Filename
		out, _ := os.Create(tmpPath)
		io.Copy(out, file)
		out.Close()

		AddMessage(fmt.Sprintf("You sent file: %s", header.Filename))

		if connSendFile != nil {
			go connSendFile(tmpPath)
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Web UI running at http://localhost" + addr)
	http.ListenAndServe(addr, nil)
}
