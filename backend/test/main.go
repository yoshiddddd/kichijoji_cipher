package main


import (
    "log"
    "net/http"
	 "encoding/json"
	 "fmt"
	  "bytes"
)
type AnswerMessage struct {
	ClientId string `json:"clientId"`
	Answer string `json:"answer"`
	Keyword string `json:"keyword"`
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
func sendToDify(data []AnswerMessage) error {
	token :="app-2FxWnRThx5ju4Wd6kMiXAIjd"
	query := fmt.Sprintf("keyword: %s client(%s) Answer: %s client(%s) Answer: %s",
	data[0].Keyword, data[0].ClientId, data[0].Answer, data[1].ClientId, data[1].Answer)
	payload := DifyRequestPayload{
        Inputs:         map[string]interface{}{}, // 必要ならば `inputs` の詳細を設定
        Query:          query,
        ResponseMode:   "blocking",
        ConversationID: "",         // 必要に応じて設定
        User:           "kyoshida",  // 必要に応じて設定
        Files: []File{
            {
                Type:           "image",
                TransferMethod: "remote_url",
                URL:            "https://cloud.dify.ai/logo/logo-site.png",
            },
        },
    }

	requestBody, err := json.Marshal(payload)
	// requestBody, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("error encoding data to JSON: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.dify.ai/v1/chat-messages", bytes.NewBuffer(requestBody))
    if err != nil {
        return fmt.Errorf("error creating HTTP request: %v", err)
    }

    // Authorizationヘッダーにトークンを設定
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Content-Type", "application/json")

    // リクエストの送信
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request to Dify: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send data to Dify, status code: %d", resp.StatusCode)
    }

    log.Println("Successfully sent data to Dify")
    log.Println(resp.Body)
    return nil
}

func main() {
	// テスト用のデータを作成
	answers := []AnswerMessage{
		{
			Keyword:  "test",
			ClientId: "client-1",
			Answer:   "answer-1",
		},
		{
			Keyword:  "test",
			ClientId: "client-2",
			Answer:   "answer-2",
		},
	}

	// テスト用のデータを送信
	err := sendToDify(answers)
	if err != nil {
		log.Printf("Error sending data to Dify: %v", err)
	}
}