package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Messages struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Max_tokens int       `json:"max_tokens"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	messages := Messages{
		Model:      "gpt-4",
		Messages:   []Message{},
		Max_tokens: 200}

	openai_api_key, exists := os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		fmt.Println("OPENAI_API_KEY environment variable does not exist. Please set it through `export OPENAI_API_KEY=sk-****`")
		os.Exit(1)
	}
	for {
		fmt.Print("User: ")

		sentence, _ := reader.ReadString('\n')

		message := Message{Role: "user", Content: sentence}

		check := strings.Replace(message.Content, "\n", "", -1)

		if strings.Compare("quit", check) == 0 {
			break
		}

		messages.Messages = append(messages.Messages, message)

		resp := process(messages, openai_api_key)

		fmt.Print("Assistant: ")

		messages.Messages = append(messages.Messages, Message{Role: "assistant", Content: resp})

		fmt.Println(resp)
	}

}

func process(messages Messages, openai_api_key string) string {
	url := "https://api.openai.com/v1/chat/completions"

	json_data, err := json.Marshal(messages)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openai_api_key))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content
	}

	return ""
}
