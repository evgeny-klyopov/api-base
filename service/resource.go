package service

import (
	"api-base/database"
	"github.com/jinzhu/gorm"
)

type Resource struct {
	db *gorm.DB
}

type ResourceInterface interface {
	GetDB() *gorm.DB
}

func (r *Resource) GetDB() *gorm.DB {
	if r.db == nil {
		r.db, _ = database.Initialize()
	}

	return r.db
}