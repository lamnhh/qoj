package config

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitialiseDatabaseConnection() {
	_ = godotenv.Load()

	connectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
}