package persistence

import (
	"fmt"

	"github.com/mendelgusmao/cetesb-telegram-bot/persistence/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Persistence struct {
	*gorm.DB
	path   string
	models []any
}

func NewPersistence(path string) *Persistence {
	return &Persistence{
		path:   path,
		models: models.Models,
	}
}

func (d *Persistence) Init() error {
	db, err := gorm.Open(sqlite.Open(d.path), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("[persistence.Persistence.Init] %v", err)
	}

	for _, model := range d.models {
		if err := db.AutoMigrate(&model); err != nil {
			return fmt.Errorf("[persistence.Persistence.Init] %v", err)
		}
	}

	d.DB = db

	return nil
}
