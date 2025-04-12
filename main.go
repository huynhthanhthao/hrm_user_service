package main

import (
	"fmt"
	"log"
	"os"

	"user_service/database"
	"user_service/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Xây dựng DSN từ biến môi trường
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Khởi tạo kết nối database
	gormDB, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Khởi động server
	if err := server.StartServer("50051", gormDB); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}