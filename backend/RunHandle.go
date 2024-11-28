package main

import (
	"encoding/json"
	"log"

)

func (s *Server) handleRegister(client *Client) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.clients[client] = true
	// s.rooms[client.RoomLevel] = append(s.rooms[client.RoomLevel], client)
	s.secretWordQueues[client.RoomLevel][client.SecretWord] = append(s.secretWordQueues[client.RoomLevel][client.SecretWord], client)
    log.Printf("Client connected: %v", client.conn.RemoteAddr())
    log.Printf("Number of clients: %v", len(s.clients))
	// log.Printf("type %f", client.Type)
	if _, ok := s.answersPerRoom[client.RoomLevel]; !ok {
		s.answersPerRoom[client.RoomLevel] = make(map[string]map[*Client]AnswerMessage)
	}
	
	if _, ok := s.answersPerRoom[client.RoomLevel][client.SecretWord]; !ok {
        s.answersPerRoom[client.SecretWord] = make(map[*Client]AnswerMessage)
    }
    // 2人のクライアントが接続されたらゲーム開始
	log.Printf("len(s.rooms[client.RoomLevel]) %d", len(s.secretWordQueues[client.RoomLevel][client.SecretWord]))
    // if len(s.clients) == s.expectedAnswerCount {
	if(len(s.secretWordQueues[client.RoomLevel][client.SecretWord]) == s.expectedAnswerCount){
        log.Printf("Start game")
        s.startGame(client.RoomLevel, client.SecretWord)
    }
}

func (s *Server) startGame(RoomLevel int, SecretWord string) {
	var sendKeyword string
	if RoomLevel == 1 {
		sendKeyword = firstRandomWordGenerate()
	} else if RoomLevel == 2 {
		sendKeyword = secondRandomWordGenerate()
	} else if RoomLevel == 3 {
		sendKeyword = thirdRandomWordGenerate()
	}
    go s.sendStartMessageToClients(sendKeyword, RoomLevel)
}

func (s *Server) sendStartMessageToClients(sendKeyword string , RoomLevel int, SecretWord string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    var msg ClientSendMessage
    msg.Signal = "start"
    msg.Word = sendKeyword

    for _, client := range s.secretWordQueues[RoomLevel][SecretWord] {
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
		//TODO room増えたらここは修正必要あり
		delete(s.rooms, client.RoomLevel)
        close(client.send)
        log.Printf("Client removed: %v", client.conn.RemoteAddr())
        log.Printf("Number of clients: %v", len(s.clients))
    }
}
