package main

import (
	"coinbase-etl/repository"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Open a connection to the PostgreSQL database
	// Connect to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	 // Automatically create the tables for your models
	 err = db.AutoMigrate(&repository.Candles{})
	 if err != nil {
		 panic("failed to auto migrate database")
	 }

}