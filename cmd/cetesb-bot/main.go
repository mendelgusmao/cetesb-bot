package main

import (
	"fmt"

	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
	"github.com/mendelgusmao/scoredb/lib/database"
)

func main() {
	scraper := scraper.New()
	database := database.NewDatabase()
	store := store.New(database, scraper)

	err := store.ScrapeAndStore()

	if err != nil {
		fmt.Println(err)
	}

	store.StartUpdater()

	cities, err := database.Query("cities", "ubatuba")
	fmt.Println(err, cities)
}
