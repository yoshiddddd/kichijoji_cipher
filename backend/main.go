package main

import (
    "log"
    "net/http"
	"encoding/json"
	"os"
    // "sync"
    "github.com/gorilla/websocket"
	// "golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // 開発環境用にCORSチェックをスキップ
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
func doesStringExist(secretWordQueues map[int]map[string][]*Client, target string, level int) bool {
    // 指定された level が存在するかをチェック
    if innerMap, exists := secretWordQueues[level]; exists {
        // target が innerMap に存在し、対応するスライスの長さが 2 であるかを判定
        if clients, ok := innerMap[target]; ok && len(clients) == 2 {
            return true
        }
    }
    return false
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
	if doesStringExist(server.secretWordQueues, registerMessage.Data.SecretWord, registerMessage.Data.Level) {
		var msg ClientSendMessage
		msg.Signal = "alreadyExist"
		msg.Word = "すでに同じワードが存在します"
		msgJson, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte(msgJson))
		log.Printf("SecretWord already exists")
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
    // PORT環境変数からポートを取得。未設定時は8080を使用
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    server := NewServer()
    go server.run()

    // 静的ファイルの提供
    http.Handle("/", http.FileServer(http.Dir("static")))

    // WebSocketエンドポイント
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(server, w, r)
    })
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    // サーバー起動
    addr := ":" + port
    log.Printf("Server starting on %s", addr)
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// func main() {
//     server := NewServer()
//     go server.run()

//     // 静的ファイルの提供
//     http.Handle("/", http.FileServer(http.Dir("static")))
    
//     // WebSocketエンドポイント
//     http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//         serveWs(server, w, r)
//     })

//     // サーバー起動
//     log.Printf("Server starting on :8080")
//     err := http.ListenAndServe(":8080", nil)
//     if err != nil {
//         log.Fatal("ListenAndServe: ", err)
//     }
// }