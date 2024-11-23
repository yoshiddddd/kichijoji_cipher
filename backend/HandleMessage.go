package main

import (
	"encoding/json"
	"log"
)

func (s *Server) handleMessage(c *Client, message []byte) {
    var receivedMessage AnswerMessage
    err := json.Unmarshal(message, &receivedMessage)
    if err != nil {
        log.Printf("Error unmarshalling message: %v", err)
        return
    }

    // ロックを取得して共有リソースを操作
    s.mutex.Lock()
    defer s.mutex.Unlock()

    // 回答を追加
    s.answers = append(s.answers, receivedMessage)
    log.Printf("Received message from client %s: %s", c.conn.RemoteAddr().String(), message)
	s.answersPerRoom[c.RoomLevel][c] = receivedMessage
    // すべての回答が揃ったかチェック
    if len(s.answersPerRoom[c.RoomLevel]) >= s.expectedAnswerCount {
        // 回答が揃った場合の処理を別の関数で行う
        s.processAnswers(c)
		delete(s.answersPerRoom, c.RoomLevel)
		delete(s.rooms, c.RoomLevel)
    }
}
func (s *Server) processAnswers(c *Client) {
    // クライアントに "end" シグナルを送信
    s.broadcastToClients(ClientSendMessage{
        Signal: "end",
        Word:   "AIが答えを出力中です",
    },c)

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
    }, c)

    // 回答リストをクリア
    s.answers = nil
}
func (s *Server) broadcastToClients(message interface{}, c *Client) {
    msgJson, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshalling message: %v", err)
        return
    }

    for client := range s.answersPerRoom[c.RoomLevel] {
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