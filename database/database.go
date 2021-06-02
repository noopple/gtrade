package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=root dbname=postgres1-db sslmode=disable"
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})

	db.Debug()

	if err != nil {
		return nil, err
	}

	return db, nil
}