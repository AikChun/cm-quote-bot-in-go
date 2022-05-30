package quoteservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aikchun/cm-quote-bot-in-go/pkg/quote"
	"github.com/go-pg/pg/v10"
)

type QuoteService struct {
	DB       *pg.DB
	UserID   int64
	UserName string
}

func (qs *QuoteService) GetUserQuotes() ([]quote.Quote, error) {
	var q []quote.Quote
	err := qs.DB.Model(&q).Where("user_id = ?", qs.UserID).Select()
	return q, err
}

func (qs *QuoteService) GetUserLatestQuote() (quote.Quote, error) {
	var q quote.Quote
	err := qs.DB.Model(&q).Where("user_id = ?", qs.UserID).Order("message_sent_at DESC").Limit(1).Select()
	return q, err
}

func (qs *QuoteService) GetQuotesByText(t string) ([]quote.Quote, error) {
	var q []quote.Quote

	err := qs.DB.Model(&q).Where("user_id = ?", qs.UserID).Where("LOWER(text) LIKE ?", "%"+t+"%").Select()

	return q, err
}

func (qs *QuoteService) SaveQuote(q *quote.Quote) error {
	_, err := qs.DB.Model(q).Insert()
	return err
}

func (qs *QuoteService) GetUserID() int64 {
	return qs.UserID
}

func (qs *QuoteService) FormatQuote(text string, year int) string {
	return fmt.Sprintf("\"%s\" - %s %d", text, qs.UserName, year)
}

func (qs *QuoteService) PickQuoteAtRandom(q []quote.Quote) quote.Quote {
	rand.Seed(time.Now().Unix())
	return q[rand.Intn(len(q))]
}
