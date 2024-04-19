package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectDB() *sql.DB {
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	DB_MASTER_PASSWORD := os.Getenv("DB_MASTER_PASSWORD")

	var dbHost string = "localhost"
	var dbPort int = 5432
	var dbUser string = DB_USER
	var dbName string = DB_NAME
	var dbPassword string = DB_MASTER_PASSWORD

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	fmt.Println(db)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
