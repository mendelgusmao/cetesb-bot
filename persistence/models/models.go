package models

import "gorm.io/gorm"

var Models = []any{
	Beach{},
}

type Beach struct {
	gorm.Model
	City    string
	Name    string
	Place   string
	Quality string
}
