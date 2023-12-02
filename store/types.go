package store

import (
	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/scoredb/lib/database"
	"github.com/mendelgusmao/scoredb/lib/fuzzymap/normalizer"
)

var databaseConfiguration = database.Configuration{
	UseLevenshtein: true,
	GramSizeLower:  2,
	GramSizeUpper:  3,
	MinScore:       0.5,
	SetConfiguration: normalizer.SetConfiguration{
		Synonyms: map[string]string{
			"avenida":       "av",
			"Caraguatatuba": "caragua",
		},
		StopWords: []string{
			"de",
			"da",
			"do",
			"das",
			"dos",
			"r",
			"rua",
			"m",
			"carlota",
		},
		Transliterate: true,
	},
}

type Store struct {
	database     *database.Database
	scraper      *scraper.Scraper
	lastChecksum uint32
}
