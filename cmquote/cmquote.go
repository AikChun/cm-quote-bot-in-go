package cmquote

import (
	"errors"
	"github.com/aikchun/cm-quote-bot-in-go/quotes"
	"github.com/aikchun/cm-quote-bot-in-go/telegrambot"
	"math/rand"
	"os"
	"strconv"
	"time"
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

func SaveQuote(u *telegrambot.Update) error {
	message := u.Message.ReplyToMessage
	if message == nil {
		return errors.New("Reply to a message to save a quote")
	}
	from := message.From
	text := message.Text

	cmID := os.Getenv("CM_ID")
	id, _ := strconv.ParseInt(cmID, 10, 64)

	if from.ID != id {
		return errors.New("This person is not CM")
	}

	q := quotes.Quote{
		Text:      text,
		UserID:    id,
		CreatedAt: time.Now(),
	}

	return quotes.SaveQuote(&q)
}
