package main

import (
	"os"

	"github.com/mendelgusmao/cetesb-telegram-bot/bot"
	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
	"github.com/mendelgusmao/scoredb/lib/database"
)

func main() {
	scraper := scraper.New()
	database := database.NewDatabase()
	store := store.New(database, scraper)
	bot := bot.New(store, os.Getenv("TELEGRAM_API_TOKEN"))
	done := make(chan bool)

	store.Work()
	bot.Work()

	<-done
}
