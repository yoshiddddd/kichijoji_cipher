package main

import (
    "log"
    "net/http"
    "sync"
	 "encoding/json"
	  "math/rand"
	  "time"
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

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage, 0, 2),
    }
}

//GPTに書かせた
func randomWordGenerate() string {
    words := []string{"下北沢", "ヘッドフォン", "データベース", "マックブック","スマートフォン","ノートパソコン","ワイヤレスイヤホン","ワイヤレスマウス","ワイヤレスキーボード"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func (s *Server) run() {
	var msg ClientSendMessage 
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
				sendKeyword := randomWordGenerate()
                
                // 各クライアントにメッセージを送信
                go func() {
                    for client := range s.clients {
                        // クライアントごとにClientIdを設定してメッセージをエンコード
                        msg.ClientId = client.conn.RemoteAddr().String()
						msg.Signal = "start"
						msg.Word = sendKeyword
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

var (
    count     userCount
    countLock sync.Mutex // 排他制御用のMutex
)
func (c *Client) readPump(s *Server) {
    defer func() {
		s.unregister <- c
        c.conn.Close()
    }()
	var msg ClientSendMessage
	var resultmsg ResultSendMessage
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
			for client := range s.clients {
				msg.ClientId = client.conn.RemoteAddr().String()
				msg.Signal = "end"
				msg.Word = "AIが答えを出力中です"
				msgJson, err := json.Marshal(msg)
				if err != nil {
					log.Printf("Error marshalling message: %v", err)
					continue
				}
				client.send <- string(msgJson)
			}
            log.Printf("Game set")
			count = 0
			log.Printf("Answers: %v", s.answers)
			answer ,err := sendToDify(s.answers)
			if err != nil {
				log.Printf("Error sending data to Dify: %v", err)
			}
			log.Printf("Answer from Dify: %s", answer)
			for client := range s.clients {
				resultmsg.ClientId = client.conn.RemoteAddr().String()
				resultmsg.Signal = "result"
				resultmsg.Word = answer
				msgJson, err := json.Marshal(resultmsg)
				if err != nil {
					log.Printf("Error marshalling message: %v", err)
					continue
				}
				client.send <- string(msgJson)
			}
			s.answers = nil
        }
        countLock.Unlock()
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