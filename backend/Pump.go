package main
import(
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
)
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
	var receivedMsg UserJoinMessage
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
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}
		log.Printf("こんにちは %s: %s", c.conn.RemoteAddr().String(), receivedMsg.Data.Name)
        // 受信したメッセージを処理する関数を呼び出す
		if(receivedMsg.Type == "answer"){
        go s.handleMessage(c, message)
		}
    }
}