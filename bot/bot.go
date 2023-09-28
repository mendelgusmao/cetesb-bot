package bot

import (
	"github.com/mendelgusmao/scoredb/lib/database"
)

func New(database *database.Database) *Bot {
	return &Bot{database: database}
}
