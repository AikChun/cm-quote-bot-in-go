package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Text string `json:"text"`
	Chat `json:"chat"`
}

type Chat struct {
	Id int64 `json:"id"`
}

type Event struct {
	Message `json:"message"`
}

type Response struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Command struct {
	FuncName string
	Args     []string
}

func HandleRequest(ctx context.Context, e Event) {
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	sendMessageUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BOT_TOKEN)

	response := Response{
		ChatId: e.Message.Chat.Id,
		Text:   e.Message.Text,
	}

	responseByteArray, err := json.Marshal(response)

	resp, err := http.Post(sendMessageUrl, "application/json", bytes.NewBuffer(responseByteArray))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

}

func SendMessage(body io.Reader) (resp *http.Response, err error) {
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	sendMessageUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BOT_TOKEN)
	return http.Post(sendMessageUrl, "application/json", body)
}

func GetFuncName(s string) (string, error) {
	var e error
	e = nil
	f := string(s[0])
	if f != "/" {
		e = errors.New("Function name don't start with '/'")
		return "", e
	}

	return string(s[1:]), e

}

func ParseCommand(text string) (Command, error) {
	var c Command
	t := strings.Trim(text, " ")
	tokens := strings.Split(t, " ")
	rawFuncName := tokens[0]

	funcName, err := GetFuncName(rawFuncName)

	if err != nil {
		return c, err
	}

	c = Command{
		FuncName: funcName,
	}

	return c, nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	var e Event

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}

	if err := json.Unmarshal(body, &e); err != nil {
		panic(err)
		return
	}

	response := Response{
		ChatId: e.Message.Chat.Id,
	}

	t := e.Message.Text
	c, err := ParseCommand(t)

	if err != nil {
		response.Text = "Enter a command"
	} else {
		response.Text = fmt.Sprintf("You did %s", c.FuncName)
	}

	responseByteArray, err := json.Marshal(response)

	resp, err := SendMessage(bytes.NewBuffer(responseByteArray))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Couldn't find .env")
	}

	e, ok := os.LookupEnv("ENVIRONMENT")

	if !ok {
		e = "dev"
	}

	if e != "dev" {
		fmt.Printf("lambda")
		lambda.Start(HandleRequest)
	} else {
		http.HandleFunc("/bot", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
