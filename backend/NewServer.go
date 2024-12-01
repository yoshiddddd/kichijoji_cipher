package main


func NewServer() *Server {
	server:= &Server{
		clients:    make(map[*Client]bool),
		rooms:      make(map[int][]*Client),
        broadcast:  make(chan string),
        register:   make(chan *Client),
        unregister: make(chan *Client),
		answers: make([]AnswerMessage,0),
		answersPerRoom: make(map[int]map[string]map[*Client]AnswerMessage),
		secretWordQueues: make(map[int]map[string][]*Client),
		expectedAnswerCount: 2,
    }
	server.answersPerRoom[BEGINNER] = make(map[string]map[*Client]AnswerMessage)
	server.answersPerRoom[INTERMEDIATE] = make(map[string]map[*Client]AnswerMessage)
	server.answersPerRoom[ADVANCED] = make(map[string]map[*Client]AnswerMessage)
	server.secretWordQueues[BEGINNER] = make(map[string][]*Client)
	server.secretWordQueues[INTERMEDIATE] = make(map[string][]*Client)
	server.secretWordQueues[ADVANCED] = make(map[string][]*Client)
	return server
}