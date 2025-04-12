package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    string `gorm:"type:uuid;primaryKey"`
	Name  string
	Email string `gorm:"unique"`
}