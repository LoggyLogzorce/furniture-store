package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const connect = "host=postgres " +
	"port=5432 " +
	"user=postgres " +
	"password=1234 " +
	"dbname=postgres " +
	"sslmode=disable"

const connectToTest = "host=postgres " +
	"port=5432 " +
	"user=postgres " +
	"password=1234 " +
	"dbname=test_db " +
	"sslmode=disable"

var database *gorm.DB
var migrate = make([]func(), 0)

func Add(mF func()) {
	migrate = append(migrate, mF)
}

func DB() *gorm.DB {
	return database
}

func Connect() {
	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	database = db
}

func Migrate() {
	for _, f := range migrate {
		f()
	}
}

func ConnectToTestDB() {
	db, err := gorm.Open(postgres.Open(connectToTest), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	database = db
}
