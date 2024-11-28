package main

import (
    "log"
    "net/http"
	"encoding/json"
    // "sync"
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

func serveWs(server *Server, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
	var registerMessage UserJoinMessage
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error reading message: %v", err)
	}
	err = json.Unmarshal(message, &registerMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}
	log.Printf("こんにちは %s: %s", conn.RemoteAddr().String(), registerMessage.Data.Level)
    client := &Client{
        conn: conn,
        send: make(chan string, 256), // バッファ付きチャネル
		RoomLevel: registerMessage.Data.Level,
		SecretWord:registerMessage.Data.SecretWord,
    }

	//ここに登録された時点でrun関数のhandleRegisterが呼ばれる
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