package main


func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		rooms:      make(map[int][]*Client),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage, 0, 2),
		answersPerRoom: make(map[int]map[*Client]AnswerMessage),
		expectedAnswerCount: 2,
    }
}