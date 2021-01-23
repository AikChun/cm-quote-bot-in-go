package telegrambot

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Bot struct {
	BaseUrl  string
	Token    string
	FuncName string
	Args     []string
}

func NewBot(token string) Bot {
	return Bot{
		Token: token,
	}
}

func callAPI(bot *Bot, method string, body io.Reader) (resp *http.Response, err error) {
	URL_PATTERN := "https://api.telegram.org/bot%s/%s"
	return http.Post(fmt.Sprintf(URL_PATTERN, bot.Token, method), "application/json", body)
}

func (bot *Bot) SendMessage(body io.Reader) {
	resp, err := callAPI(bot, "sendMessage", body)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
