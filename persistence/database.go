package persistence

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Persistence struct {
	path     string
	models   []any
	database *gorm.DB
}

func NewPersistence(path string) *Persistence {
	return &Persistence{
		path:   path,
		models: models,
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

	d.database = db

	return nil
}

func (d *Persistence) CreateOrUpdateBeach(beach Beach) error {
	var tmpBeach Beach

	d.database.First(&tmpBeach, "city = ? AND name = ?", beach.City, beach.Name)
}
