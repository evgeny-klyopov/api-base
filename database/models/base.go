package models

import (
	"github.com/jinzhu/gorm"
)

type Base struct {
	gorm.Model
	Created   User   `gorm:"foreignkey:UserID"`
	CreatedBy uint `gorm:"not null;"`
	Updated   User   `gorm:"foreignkey:UserID"`
	UpdatedBy *uint
	Deleted   User   `gorm:"foreignkey:UserID"`
	DeletedBy *uint
	IsActive  bool `gorm:"not null;"`
}