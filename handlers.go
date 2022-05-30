package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	telegrambot "github.com/AikChun/yagotb"
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
