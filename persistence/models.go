package persistence

import "gorm.io/gorm"

var models = []any{
	Beach{},
}

type Beach struct {
	gorm.Model
	City    string
	Name    string
	Place   string
	Quality string
}
