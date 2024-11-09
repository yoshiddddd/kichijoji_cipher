package main

func (s *Server) run() {
    for {
        select {
        case client := <-s.register:
            s.handleRegister(client)

        case client := <-s.unregister:
            s.handleUnregister(client)

        case message := <-s.broadcast:
            s.handleBroadcast(message)
        }
    }
}