package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitCockroach() {
	connStr := os.Getenv("CONN_STR")

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Error connecting to CockroachDB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Error pinging CockroachDB:", err)
	}

	log.Println("✅ Connected to CockroachDB Cloud")
}
