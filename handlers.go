package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aikchun/cm-quote-bot-in-go/cmquote"
	"github.com/aikchun/cm-quote-bot-in-go/telegrambot"
	"log"
	"strings"
)

func Echo(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   strings.Join(args, " "),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func RandomQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
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

func LatestQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
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

func SaveQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		ReplyToMessageID: u.Message.MessageID,
	}

	if err := cmquote.SaveQuote(u); err != nil {
		response.Text = err.Error()

		responseByteArray, _ := json.Marshal(response)

		bot.SendMessage(bytes.NewBuffer(responseByteArray))
		return
	}

	response.Text = "Saved!"

	responseByteArray, _ := json.Marshal(response)

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func Help(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
	helpText := fmt.Sprintf("/randomquote - Get an awesome quote!\n/latestquote - Get the latest addition!\n/savequote - Save a CM Quote\n/help - Get a list of commands")
	payload := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   helpText,
	}

	responseByteArray, _ := json.Marshal(payload)

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}
