# socket-chat

A lightweight, two-way TCP chat application in Go with both a CLI and a web UI interface, and support for file transfers.

---

## Features

- **Text messaging over TCP** between two clients.
- **Dual interface**:  
  - Terminal (CLI) for interactive chat.  
  - Browser-based UI with HTML, JS, and HTTP polling—no WebSockets required.
- **File transfer support**: Send and receive files via both CLI and Web UI.

---

## Overview

`socket-chat` allows two peers to communicate over TCP using a simple text-based protocol ever-where the peers can exchange messages and files. The app supports both CLI mode—prompt-based human chats—and browser-based UI for a modern experience, with file upload and automatic download links.

---

## Prerequisites

- Go (version 1.18 or newer recommended)  
- A modern web browser to access the Web UI  
- Terminal access for CLI mode

---

## Setup & Installation

1. Clone your repository:
    ```bash
   git clone https://github.com/SaeedMPro/socket-chat.git
   cd socket-chat
    ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```
3. Prepare configuration by editing `config/config.json`. Example:

   ```json
   {
     "client-one": {
       "host": "localhost",
       "port": 8000,
       "ui_port": 8080
     },
     "client-two": {
       "host": "localhost",
       "port": 9000,
       "ui_port": 9090
     }
   }
   ```

---

## Usage

Run two clients in separate terminal windows:

```bash
go run main.go client-one
```

In another:

```bash
go run main.go client-two
```

* The CLI will open at ports 8000 and 9000 respectively.
* The Web UI will serve on ports 8080 and 9090, e.g. `http://localhost:8080`.

Open the URLs in your browser to use the web interface.

---

## Web UI Features

* **Send messages**: Type your message and click **Send**—appears in both your browser and the peer's screen.
* **Send files**: Use the file picker to upload a file, which streams directly to the peer.
* **Auto-update messages**: Polls the server every second to fetch new chat messages seamlessly.
* **Responsive & styled UI**: Includes message bubbles, clean layout, and file upload support.

---

## Project Structure

```
socket-chat/
├── config/
│   ├── config.go
│   └── config.json       # Client configurations (ports, etc.)
├── internal/
│   ├── client/
│   │   ├── client.go     # Core chat logic (listener, sender)
│   │   └── connection.go # Optional details
│   ├── transport/
│   │   ├── file.go       # File sending/receiving utilities
│   │   └── message.go    # Text message utilities
│   └── ui/
│       ├── cli.go        # Terminal UI logic
│       ├── web.go        # HTTP UI server with message/file handlers
│       └── index.html    # Browser interface
├── main.go
├── go.mod
└── README.md
```

