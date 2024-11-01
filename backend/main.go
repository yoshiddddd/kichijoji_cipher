package main

import (
    "log"
    "net/http"
    "sync"
	 "encoding/json"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // 開発環境用にCORSチェックをスキップ
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Client struct {
    conn *websocket.Conn
    send chan string
}

type Server struct {
    clients    map[*Client]bool
    broadcast  chan string
    register   chan *Client
    unregister chan *Client
    mutex      sync.Mutex
}
type Message struct {
    Signal string `json:"signal"`
    Word   string `json:"word"`
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (s *Server) run() {
	msg := Message{
		Signal: "start",
		Word:   "爽健美茶",
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

    for {
        select {
        case client := <-s.register:
            s.mutex.Lock()
            s.clients[client] = true
            log.Printf("Client connected: %v", client.conn.RemoteAddr())
            log.Printf("Number of clients: %v", len(s.clients))
            if len(s.clients) == 2 {
                log.Printf("start game")
                // ブロードキャストメッセージを送信
                go func() {
                    for client := range s.clients {
                        select {
                        // case client.send <- "start":
                        case client.send <- string(msgJson):
                            log.Printf("Message sent to client: %v", client.conn.RemoteAddr())
                        default:
                            s.mutex.Lock()
                            delete(s.clients, client)
                            close(client.send)
                            s.mutex.Unlock()
                            log.Printf("Failed to send message to client: %v", client.conn.RemoteAddr())
                        }
                    }
                }()
            }
            s.mutex.Unlock()

        case client := <-s.unregister:
            s.mutex.Lock()
            if _, ok := s.clients[client]; ok {
                delete(s.clients, client)
                close(client.send)
                log.Printf("Client disconnected: %v", client.conn.RemoteAddr())
                log.Printf("Number of clients: %v", len(s.clients))
            }
            s.mutex.Unlock()

        case message := <-s.broadcast:
            log.Printf("Broadcasting message: %v", message)
            s.mutex.Lock()
            for client := range s.clients {
                select {
                case client.send <- message:
                    log.Printf("Broadcast message sent to client: %v", client.conn.RemoteAddr())
                default:
                    delete(s.clients, client)
                    close(client.send)
                    log.Printf("Failed to broadcast message to client: %v", client.conn.RemoteAddr())
                }
            }
            s.mutex.Unlock()
        }
    }
}

func (c *Client) writePump() {
    defer func() {
        c.conn.Close()
    }()

    for {
        message, ok := <-c.send
        if !ok {
            // チャネルが閉じられている
            c.conn.WriteMessage(websocket.CloseMessage, []byte{})
            return
        }

        err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
        if err != nil {
            log.Printf("Error writing message: %v", err)
            return
        }
    }
}

func (c *Client) readPump(s *Server) {
    defer func() {
        s.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
            }
            break
        }
        log.Printf("Received message from client: %s", message)
        // 必要に応じてメッセージを処理
    }
}

func serveWs(server *Server, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }

    client := &Client{
        conn: conn,
        send: make(chan string, 256), // バッファ付きチャネル
    }

    server.register <- client

    // クライアントの送受信を開始
    go client.writePump()
    go client.readPump(server)
}

func main() {
    server := NewServer()
    go server.run()

    // 静的ファイルの提供
    http.Handle("/", http.FileServer(http.Dir("static")))
    
    // WebSocketエンドポイント
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(server, w, r)
    })

    // サーバー起動
    log.Printf("Server starting on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}