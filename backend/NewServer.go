package main


func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage, 0, 2),
		expectedAnswerCount: 2,
    }
}