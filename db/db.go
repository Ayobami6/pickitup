package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb(dbStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbStr), &gorm.Config{})
	if err != nil {
        return nil, err
    }
	log.Println("Database Connected Successfully!")
	return db, nil

}