package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate để tạo bảng
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}