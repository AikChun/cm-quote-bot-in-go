package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	telegrambot "github.com/AikChun/yagotb"
	"github.com/aikchun/cm-quote-bot-in-go/pkg/quotecontroller"
	"github.com/aikchun/cm-quote-bot-in-go/pkg/quoteservice"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

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
	bot, err := telegrambot.NewBot(os.Getenv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	quotedUserID := os.Getenv("QUOTED_USER_ID")

	userID, err := strconv.ParseInt(quotedUserID, 10, 64)

	if err != nil {
		panic(err)
	}

	userName := os.Getenv("QUOTED_USER_NAME")

	url := os.Getenv("DATABASE_URL")
	opt, err := pg.ParseURL(url)

	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)

	cmQuoteController := quotecontroller.QuoteController{
		QS: &quoteservice.QuoteService{
			DB:       db,
			UserID:   userID,
			UserName: userName,
		},
	}

	bot.AddHandler("/echo", Echo)
	bot.AddHandler("/randomquote", cmQuoteController.RandomQuote)
	bot.AddHandler("/savequote", cmQuoteController.SaveQuote)
	bot.AddHandler("/latestquote", cmQuoteController.LatestQuote)
	bot.AddHandler("/crisis", cmQuoteController.RandomCrisisQuote)
	bot.AddHandler("/help", Help)

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
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
