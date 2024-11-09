package main
import(
	"github.com/gorilla/websocket"
	"log"
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

        // 受信したメッセージを処理する関数を呼び出す
        go s.handleMessage(c, message)
    }
}