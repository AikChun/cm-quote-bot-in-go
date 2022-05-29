package cmquote

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"

	telegrambot "github.com/AikChun/yagotb"
	"github.com/aikchun/cm-quote-bot-in-go/quotes"
)

func PickQuoteAtRandom(q []quotes.Quote) quotes.Quote {
	rand.Seed(time.Now().Unix())
	return q[rand.Intn(len(q))]
}

func GetAllCMQuotes() []quotes.Quote {
	var q []quotes.Quote
	cmID := os.Getenv("CM_ID")
	id, err := strconv.ParseInt(cmID, 10, 64)
	if err != nil {
		panic(err)
	}

	if err := quotes.GetUserQuotes(&q, id); err != nil {
		panic(err)
	}

	return q

}

func GetRandomQuote() quotes.Quote {
	q := GetAllCMQuotes()
	return PickQuoteAtRandom(q)
}

func GetLatestQuote() quotes.Quote {
	var q quotes.Quote
	cmID := os.Getenv("CM_ID")
	id, err := strconv.ParseInt(cmID, 10, 64)
	if err != nil {
		panic(err)
	}

	if err := quotes.GetUserLatestQuote(&q, id); err != nil {
		panic(err)
	}

	return q
}

func GetRandomCrisisQuote() quotes.Quote {
	var q []quotes.Quote
	cmID := os.Getenv("CM_ID")
	id, err := strconv.ParseInt(cmID, 10, 64)
	if err != nil {
		panic(err)
	}

	if err := quotes.GetQuotesByText(&q, id, "crisis"); err != nil {
		panic(err)
	}

	return PickQuoteAtRandom(q)

}

func SaveQuote(u *telegrambot.Update) error {
	message := u.Message.ReplyToMessage
	if message == nil {
		return errors.New("reply to a message to save a quote")
	}
	from := message.From
	text := message.Text
	date := time.Unix(message.Date, 0)

	cmID := os.Getenv("CM_ID")
	id, _ := strconv.ParseInt(cmID, 10, 64)

	if from.ID != id {
		return errors.New("this person is not CM")
	}

	q := quotes.Quote{
		Text:          text,
		UserID:        id,
		CreatedAt:     time.Now(),
		MessageSentAt: date,
	}

	return quotes.SaveQuote(&q)
}
