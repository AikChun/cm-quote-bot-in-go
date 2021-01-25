package telegrambot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Bot struct {
	Token    string
	Username string
	Handlers map[string]func(*Bot, *Update, []string)
}

func NewBot(token string, username string) Bot {
	return Bot{
		Token:    token,
		Username: username,
		Handlers: make(map[string]func(*Bot, *Update, []string)),
	}
}

type Update struct {
	UpdateID          int64   `json:"update_id"`
	Message           Message `json:"message,omitempty"`
	EditedMessage     Message `json:"edited_message,omitempty"`
	ChannelPost       Message `json:"channel_post,omitempty"`
	EditedChannelPost Message `json:"edited_channel_post,omitempty"`
}

type User struct {
	ID        int64 `json:id`
	FirstName int64 `json:first_name`
}

type Message struct {
	MessageID      int64  `json:"message_id"`
	Text           string `json:"text"`
	Chat           `json:"chat"`
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	From           User     `json:"from"`
}

type Chat struct {
	Id int64 `json:"id"`
}

type Response struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
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

func (bot *Bot) AddHandler(s string, f func(*Bot, *Update, []string)) {
	bot.Handlers[s] = f
	ss := fmt.Sprintf("%s%s", s, bot.Username)
	bot.Handlers[ss] = f
}

func (bot *Bot) HandleUpdate(u *Update) {
	s := u.Message.Text

	trimmed := strings.Trim(s, " ")
	tokens := strings.Split(trimmed, " ")
	funcName := tokens[0]
	args := tokens[1:]

	if f, ok := bot.Handlers[funcName]; ok {
		f(bot, u, args)
	}

}
