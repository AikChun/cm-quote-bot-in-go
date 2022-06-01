package quotecontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	telegrambot "github.com/AikChun/yagotb"
	"github.com/aikchun/cm-quote-bot-in-go/pkg/quote"
)

type QuoteService interface {
	GetUserID() int64
	GetUserQuotes() ([]quote.Quote, error)
	GetUserLatestQuote() (quote.Quote, error)
	GetQuotesByText(t string) ([]quote.Quote, error)
	SaveQuote(q *quote.Quote) error
	FormatQuote(string, int) string
	PickQuoteAtRandom(q []quote.Quote) quote.Quote
}

type QuoteController struct {
	QS QuoteService
}

func (qc *QuoteController) LatestQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {

	q, err := qc.QS.GetUserLatestQuote()

	if err != nil {
		log.Fatal(err)
	}

	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   qc.QS.FormatQuote(q.Text, q.MessageSentAt.Year()),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func (qc *QuoteController) RandomQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {

	quotes, err := qc.QS.GetUserQuotes()

	if err != nil {
		log.Fatal(err)
	}

	q := qc.QS.PickQuoteAtRandom(quotes)

	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   qc.QS.FormatQuote(q.Text, q.MessageSentAt.Year()),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))
}

func (qc *QuoteController) RandomCrisisQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {

	quotes, err := qc.QS.GetQuotesByText("crisis")

	if err != nil {
		log.Fatal(err)
	}

	q := qc.QS.PickQuoteAtRandom(quotes)

	response := telegrambot.SendMessagePayload{
		ChatId: u.Message.Chat.Id,
		Text:   qc.QS.FormatQuote(q.Text, q.MessageSentAt.Year()),
	}

	responseByteArray, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}

func (qc *QuoteController) SaveQuote(bot *telegrambot.Bot, u *telegrambot.Update, args []string) {
	var err error

	message := u.Message.ReplyToMessage
	if message == nil {
		err = errors.New("reply to a message to save a quote")
	}

	response := telegrambot.SendMessagePayload{
		ChatId:           u.Message.Chat.Id,
		ReplyToMessageID: u.Message.MessageID,
	}

	from := message.From
	text := message.Text
	date := time.Unix(message.Date, 0)

	id := qc.QS.GetUserID()

	if from.ID != id {
		response.Text = "this person is not CM"

		responseByteArray, _ := json.Marshal(response)

		bot.SendMessage(bytes.NewBuffer(responseByteArray))
		return
	}

	q := &quote.Quote{
		Text:          text,
		UserID:        id,
		CreatedAt:     time.Now(),
		MessageSentAt: date,
	}

	err = qc.QS.SaveQuote(q)

	if err != nil {
		response.Text = err.Error()

		responseByteArray, _ := json.Marshal(response)

		bot.SendMessage(bytes.NewBuffer(responseByteArray))
		return
	}

	response.Text = "Saved!"

	responseByteArray, _ := json.Marshal(response)

	bot.SendMessage(bytes.NewBuffer(responseByteArray))

}
