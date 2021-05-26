package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	telegrambot "github.com/AikChun/yagotb"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func CreateTelegramBot() *telegrambot.Bot {
	bot, err := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	bot.AddHandler("/echo", Echo)
	bot.AddHandler("/randomquote", RandomQuote)
	bot.AddHandler("/latestquote", LatestQuote)
	bot.AddHandler("/savequote", SaveQuote)
	bot.AddHandler("/help", Help)

	return bot
}

func HandleRequest(ctx context.Context, u telegrambot.Update) {
	HandleTelegramUpdate(&u)
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	HandleTelegramUpdate(&u)
}

func HandleTelegramUpdate(u *telegrambot.Update) {
	bot := CreateTelegramBot()
	bot.HandleUpdate(u)
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
		lambda.Start(HandleRequest)
	} else {
		http.HandleFunc("/bot", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
