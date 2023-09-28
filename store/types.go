package store

import (
	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/scoredb/lib/database"
	"github.com/mendelgusmao/scoredb/lib/fuzzymap/normalizer"
)

var databaseConfiguration = database.Configuration{
	UseLevenshtein: true,
	GramSizeLower:  3,
	GramSizeUpper:  5,
	MinScore:       50,
	SetConfiguration: normalizer.SetConfiguration{
		Synonyms: map[string]string{
			"avenida": "av",
		},
		StopWords:     []string{"de"},
		Transliterate: true,
	},
}

type Store struct {
	database *database.Database
	scraper  *scraper.Scraper
}
