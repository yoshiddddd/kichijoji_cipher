package main

import (
	"encoding/json"
	"log"

)

func (s *Server) handleRegister(client *Client) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.clients[client] = true
    log.Printf("Client connected: %v", client.conn.RemoteAddr())
    log.Printf("Number of clients: %v", len(s.clients))

    // 2人のクライアントが接続されたらゲーム開始
    if len(s.clients) == 2 {
        log.Printf("Start game")
        s.startGame()
    }
}

func (s *Server) startGame() {
    sendKeyword := randomWordGenerate()
    go s.sendStartMessageToClients(sendKeyword)
}

func (s *Server) sendStartMessageToClients(sendKeyword string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    var msg ClientSendMessage
    msg.Signal = "start"
    msg.Word = sendKeyword

    for client := range s.clients {
        // クライアントごとに ClientId を設定
        msg.ClientId = client.conn.RemoteAddr().String()
        msgJson, err := json.Marshal(msg)
        if err != nil {
            log.Printf("Error marshalling message: %v", err)
            continue
        }

        // メッセージをクライアントに送信
        s.sendMessageToClient(client, string(msgJson))
    }
}
func (s *Server) sendMessageToClient(client *Client, message string) {
    select {
    case client.send <- message:
        log.Printf("Message sent to client: %v", client.conn.RemoteAddr())
    default:
        s.removeClient(client)
        log.Printf("Failed to send message to client: %v", client.conn.RemoteAddr())
    }
}
func (s *Server) handleUnregister(client *Client) {
    s.removeClient(client)
    log.Printf("Client disconnected: %v", client.conn.RemoteAddr())
}


func (s *Server) handleBroadcast(message string) {
    log.Printf("Broadcasting message: %v", message)
    s.mutex.Lock()
    defer s.mutex.Unlock()

    for client := range s.clients {
        s.sendMessageToClient(client, message)
    }
}


func (s *Server) removeClient(client *Client) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, ok := s.clients[client]; ok {
        delete(s.clients, client)
        close(client.send)
        log.Printf("Client removed: %v", client.conn.RemoteAddr())
        log.Printf("Number of clients: %v", len(s.clients))
    }
}
