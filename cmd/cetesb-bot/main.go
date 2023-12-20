package main

import (
	"log"
	"os"

	"github.com/mendelgusmao/cetesb-telegram-bot/bot"
	"github.com/mendelgusmao/cetesb-telegram-bot/persistence"
	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
	"github.com/mendelgusmao/scoredb/lib/database"
)

func main() {
	persistence := persistence.NewPersistence(os.Getenv("PERSISTENCE_PATH"))

	if err := persistence.Init(); err != nil {
		log.Fatal(err)
	}

	scraper := scraper.New()
	memory := database.NewDatabase()
	store := store.New(persistence, memory, scraper)
	bot := bot.New(store, os.Getenv("TELEGRAM_API_TOKEN"))
	done := make(chan bool)

	store.Work()
	bot.Work()

	<-done
}
