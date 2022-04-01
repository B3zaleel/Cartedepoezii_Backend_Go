package utils

import (
	"os"
	"log"
	_ "database/sql"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

// Creates a database connection.
func GetDBConnection() (db *sqlx.DB, err error) {
	db, err = sqlx.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

// Initializes the database.
func InitDB() (err error) {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	file_bytes, err := os.ReadFile("src/db/DBInit.sql")
	if err != nil {
		return err
	}
	db.Exec(string(file_bytes))
	return nil
}
