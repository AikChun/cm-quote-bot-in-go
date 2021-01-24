package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aikchun/cm-quote-bot-in-go/cmquote"
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

func HandleEvent(bot *telegrambot.Bot, u telegrambot.Update) {
	switch bot.FuncName {
	case "/echo":
		Echo(bot, u)
	case "/randomquote":
		RandomQuote(bot, u)
	case "/randomquote@cmquotebot":
		RandomQuote(bot, u)
	case "/latestquote":
		LatestQuote(bot, u)
	case "/latestquote@cmquotebot":
		LatestQuote(bot, u)
	case "/savequote":
		SaveQuote(bot, u)
	case "/savequote@cmquotebot":
		SaveQuote(bot, u)
	}

}

func Echo(bot *telegrambot.Bot, u telegrambot.Update) {
	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   strings.Join(bot.Args, " "),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func RandomQuote(bot *telegrambot.Bot, u telegrambot.Update) {
	quote := cmquote.GetRandomQuote()

	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		Text:             quote.Text,
		ReplyToMessageID: u.Message.MessageID,
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func LatestQuote(bot *telegrambot.Bot, u telegrambot.Update) {
	quote := cmquote.GetLatestQuote()

	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		Text:             quote.Text,
		ReplyToMessageID: u.Message.MessageID,
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func SaveQuote(bot *telegrambot.Bot, u telegrambot.Update) {
	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
	}

	if err := cmquote.SaveQuote(&u); err != nil {
		response.Text = err.Error()

		responseByteArray, _ := json.Marshal(response)

		bot.SendMessage(bytes.NewBuffer(responseByteArray))
		return
	}

	response.Text = "Saved!"

	responseByteArray, _ := json.Marshal(response)

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func HandleRequest(ctx context.Context, u telegrambot.Update) {
	bot := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))

	trimmed := strings.Trim(u.Message.Text, " ")
	tokens := strings.Split(trimmed, " ")
	bot.FuncName = tokens[0]
	bot.Args = tokens[1:]

	HandleEvent(&bot, u)
}

func handler(w http.ResponseWriter, r *http.Request) {
	bot := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))
	var u telegrambot.Update

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &u); err != nil {
		panic(err)
	}

	trimmed := strings.Trim(u.Message.Text, " ")
	tokens := strings.Split(trimmed, " ")
	bot.FuncName = tokens[0]
	bot.Args = tokens[1:]

	HandleEvent(&bot, u)

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
