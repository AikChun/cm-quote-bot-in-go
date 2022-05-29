package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	telegrambot "github.com/AikChun/yagotb"
	"github.com/aikchun/cm-quote-bot-in-go/cmquote"
)

func FormatQuote(text string, year int) string {
	return fmt.Sprintf("\"%s\" - CM %d", text, year)

}

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

	var targetMessageID int64

	if u.Message.ReplyToMessage != nil {
		targetMessageID = u.Message.ReplyToMessage.MessageID
	} else {
		targetMessageID = u.Message.MessageID
	}

	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		Text:             FormatQuote(quote.Text, quote.MessageSentAt.Year()),
		ReplyToMessageID: targetMessageID,
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
		Text:             FormatQuote(quote.Text, quote.MessageSentAt.Year()),
		ReplyToMessageID: u.Message.MessageID,
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func RandomCrisisQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
	quote := cmquote.GetRandomCrisisQuote()

	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		Text:             FormatQuote(quote.Text, quote.MessageSentAt.Year()),
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
	helpText := fmt.Sprintf("" +
		"/randomquote - Get an awesome quote!\n" +
		"/savequote - Save a CM Quote\n" +
		"/latestquote - Get the latest addition!\n" +
		"/crisis - CRISIS!!!\n" +
		"/help - Get a list of commands")
	payload := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   helpText,
	}

	responseByteArray, _ := json.Marshal(payload)

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}
