package main

import (
    "log"
    "net/http"
    // "sync"
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

func (s *Server) handleMessage(c *Client, message []byte) {
    var sendMessage AnswerMessage
    err := json.Unmarshal(message, &sendMessage)
    if err != nil {
        log.Printf("Error unmarshalling message: %v", err)
        return
    }

    // ロックを取得して共有リソースを操作
    s.mutex.Lock()
    defer s.mutex.Unlock()

    // 回答を追加
    s.answers = append(s.answers, sendMessage)
    log.Printf("Received message from client %s: %s", c.conn.RemoteAddr().String(), message)

    // すべての回答が揃ったかチェック
    if len(s.answers) >= s.expectedAnswerCount {
        // 回答が揃った場合の処理を別の関数で行う
        s.processAnswers()
    }
}
func (s *Server) processAnswers() {
    // クライアントに "end" シグナルを送信
    s.broadcastToClients(ClientSendMessage{
        Signal: "end",
        Word:   "AIが答えを出力中です",
    })

    log.Printf("Game set")
    log.Printf("Answers: %v", s.answers)

    // AIへのリクエストを行う
    answer, err := sendToDify(s.answers)
    if err != nil {
        log.Printf("Error sending data to Dify: %v", err)
        return
    }
    log.Printf("Answer from Dify: %s", answer)

    // クライアントに結果を送信
    s.broadcastToClients(ResultSendMessage{
        Signal: "result",
        Word:   answer,
    })

    // 回答リストをクリア
    s.answers = nil
}
func (s *Server) broadcastToClients(message interface{}) {
    msgJson, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshalling message: %v", err)
        return
    }

    for client := range s.clients {
        select {
        case client.send <- string(msgJson):
            // 送信成功
        default:
            // 送信失敗（チャネルが詰まっている場合など）
            close(client.send)
            delete(s.clients, client)
        }
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