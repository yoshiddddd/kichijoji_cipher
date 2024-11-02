package main

import (
    "log"
    "net/http"
    "sync"
	 "encoding/json"
	 "fmt"
	  "bytes"
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
	answers []AnswerMessage
}
type Message struct {
    Signal string `json:"signal"`
    Word   string `json:"word"`
}
type AnswerMessage struct {
	ClientId string `json:"clientId"`
	Answer string `json:"answer"`
	Keyword string `json:"keyword"`
}
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage, 0, 2),
    }
}
type DifyRequestPayload struct {
    Inputs         map[string]interface{} `json:"inputs"`
    Query          string                 `json:"query"`
    ResponseMode   string                 `json:"response_mode"`
    ConversationID string                 `json:"conversation_id"`
    User           string                 `json:"user"`
    Files          []File                 `json:"files"`
}

type File struct {
    Type           string `json:"type"`
    TransferMethod string `json:"transfer_method"`
    URL            string `json:"url"`
}
func (s *Server) run() {
	type Message struct {
		ClientId string `json:"clientId"`
		Signal   string `json:"signal"`
		Word     string `json:"word"`
	}
	var msg Message 
    for {
        select {
        case client := <-s.register:
            s.mutex.Lock()
            s.clients[client] = true
            log.Printf("Client connected: %v", client.conn.RemoteAddr())
            log.Printf("Number of clients: %v", len(s.clients))
            
            // 2人のクライアントが接続されたらゲーム開始
            if len(s.clients) == 2 {
                log.Printf("Start game")
                
                // 各クライアントにメッセージを送信
                go func() {
                    for client := range s.clients {
                        // クライアントごとにClientIdを設定してメッセージをエンコード
                        msg.ClientId = client.conn.RemoteAddr().String()
						msg.Signal = "start"
						msg.Word = "apple"
                        msgJson, err := json.Marshal(msg)
                        if err != nil {
                            log.Printf("Error marshalling message: %v", err)
                            continue
                        }

                        // メッセージをクライアントに送信
                        select {
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

type userCount int

var (
    count     userCount
    countLock sync.Mutex // 排他制御用のMutex
)
func (c *Client) readPump(s *Server) {
    defer func() {
		s.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
		var sendMessage AnswerMessage
		err = json.Unmarshal(message, &sendMessage)
		// log.Printf("Received message: %+v", sendMessage.ClientId)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
            }
            break
        }
        if message != nil {
            countLock.Lock()
			s.answers = append(s.answers, sendMessage)
            count++
            log.Println("Current user count:", count)
            countLock.Unlock()
        }
        log.Printf("Received message from client: %s", message)

        countLock.Lock()
        // if count == 2 {
		if len(s.answers) == 2 {
            log.Printf("Game set")
			count = 0
			log.Printf("Answers: %v", s.answers)
			err := sendToDify(s.answers)
			if err != nil {
				log.Printf("Error sending data to Dify: %v", err)
			}
        }
        countLock.Unlock()
    }
}
func sendToDify(data []AnswerMessage) error {
	token :="app-iEHIHbKk9exa7pkrPTnAdBIW"
	query := fmt.Sprintf("keyword: %s client(%s) Answer: %s client(%s) Answer: %s",
	data[0].Keyword, data[0].ClientId, data[0].Answer, data[1].ClientId, data[1].Answer)
	payload := DifyRequestPayload{
        Inputs:         map[string]interface{}{}, // 必要ならば `inputs` の詳細を設定
        Query:          query,
        ResponseMode:   "streaming",
        ConversationID: "",         // 必要に応じて設定
        User:           "abc-123",  // 必要に応じて設定
        Files: []File{
            {
                Type:           "image",
                TransferMethod: "remote_url",
                URL:            "https://cloud.dify.ai/logo/logo-site.png",
            },
        },
    }

	requestBody, err := json.Marshal(payload)
	// requestBody, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("error encoding data to JSON: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.dify.ai/v1/chat-messages", bytes.NewBuffer(requestBody))
    if err != nil {
        return fmt.Errorf("error creating HTTP request: %v", err)
    }

    // Authorizationヘッダーにトークンを設定
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Content-Type", "application/json")

    // リクエストの送信
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request to Dify: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send data to Dify, status code: %d", resp.StatusCode)
    }

    log.Println("Successfully sent data to Dify")
    return nil
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