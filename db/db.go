package db

import (
	"github.com/go-pg/pg/v10"
	"os"
)

func NewDB() *pg.DB {
	url := os.Getenv("DATABASE_URL")
	opt, err := pg.ParseURL(url)

	if err != nil {
		panic(err)
	}

	return pg.Connect(opt)
}
