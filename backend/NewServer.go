package main


func NewServer() *Server {
	server:= &Server{
		clients:    make(map[*Client]bool),
		rooms:      make(map[int][]*Client),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage, 0, 2),
		answersPerRoom: make(map[int]map[string]map[*Client]AnswerMessage)
		secretWordQueues: make(map[int]map[string][]*Client)
		expectedAnswerCount: 2,
    }
	secretWordQueues[BEGINNER] = make(map[string][]*Client)
	secretWordQueues[INTERMEDIATE] = make(map[string][]*Client)
	secretWordQueues[ADVANCED] = make(map[string][]*Client)

	return server
}