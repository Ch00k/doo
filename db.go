package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func initDB(host, port, username, password, dbname string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, username, password, dbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
}

func SetupDB() *gorm.DB {
	dbHost := GetEnvVar("DOO_DB_HOST", "localhost")
	dbPort := GetEnvVar("DOO_DB_PORT", "5432")
	dbUser := GetEnvVar("DOO_DB_USER", "doo")
	dbPass := GetEnvVar("DOO_DB_PASSWORD", "doo")
	dbName := GetEnvVar("DOO_DB_NAME", "doo")

	db, err := initDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&Entry{},
		&Comment{},
	)

	return db
}
