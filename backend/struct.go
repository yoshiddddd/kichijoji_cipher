package main
import(
	"github.com/gorilla/websocket"
	"sync"
)


type Client struct {
    conn *websocket.Conn
    send chan string
	RoomLevel int
}

type Server struct {

    clients    map[*Client]bool
	rooms      map[int][]*Client
    broadcast  chan string
    register   chan *Client
    unregister chan *Client
    mutex      sync.Mutex
	answers []AnswerMessage
	expectedAnswerCount int
}
type Message struct {
    Signal string `json:"signal"`
    Word   string `json:"word"`
}
type AnswerMessage struct {
	Type string `json:"type"`
	Data struct{
		ClientId string `json:"clientId"`
		Name string `json:"name"`
		Answer string `json:"answer"`
		Keyword string `json:"keyword"`
		CountTime int `json:"countTime"`
	} `json:"data"`
}
type UserJoinMessage struct {
	Type string `json:"type"`
	Data struct {
		Name  string `json:"name"`
		Level int    `json:"level"`
	} `json:"data"`
}

type DifyRequestPayload struct {
    Inputs         map[string]interface{} `json:"inputs"`
    Query          string                 `json:"query"`
    ResponseMode   string                 `json:"response_mode"`
    ConversationID string                 `json:"conversation_id"`
    User           string                 `json:"user"`
    Files          []File                 `json:"files"`
}

type File struct {
    Type           string `json:"type"`
    TransferMethod string `json:"transfer_method"`
    URL            string `json:"url"`
}
type ClientSendMessage struct {
	ClientId string `json:"clientId"`
	Signal   string `json:"signal"`
	Word     string `json:"word"`
}
type ResultSendMessage struct {
	ClientId string `json:"clientId"`
	Signal   string `json:"signal"`
	Word     string `json:"word"`
	Winner   string `json:"winner"`
}
type DifyResponse struct {
    Answer string `json:"answer"`
}

type userCount int