package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aikchun/cm-quote-bot-in-go/telegrambot"
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

func HandleRequest(ctx context.Context, e Event) {
	bot := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))

	response := Response{
		ChatId: e.Message.Chat.Id,
		Text:   e.Message.Text,
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func HandleEvent(bot *telegrambot.Bot, e Event) {
	switch bot.FuncName {
	case "/echo":
		Echo(bot, e)
	default:
		// freebsd, openbsd,
		// plan9, windows...
		InvalidCommand(bot, e)
	}

}

func Echo(bot *telegrambot.Bot, e Event) {
	response := telegrambot.SendMessagePayload{
		ChatId: e.Message.Chat.Id,
		Text:   strings.Join(bot.Args, " "),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func InvalidCommand(bot *telegrambot.Bot, e Event) {
	response := telegrambot.SendMessagePayload{
		ChatId: e.Message.Chat.Id,
		Text:   "Invalid command",
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func handler(w http.ResponseWriter, r *http.Request) {
	bot := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))
	var e Event

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &e); err != nil {
		panic(err)
	}

	trimmed := strings.Trim(e.Message.Text, " ")
	tokens := strings.Split(trimmed, " ")
	bot.FuncName = tokens[0]
	bot.Args = tokens[1:]

	HandleEvent(&bot, e)

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
