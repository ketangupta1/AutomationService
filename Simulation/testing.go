package Simulation

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	//serverName := "postgreSQL15"
	//username := "postgres"
	//password := "admin"
	//dbName := "HistoricalData"

	// Create the connection string
	connStr := "postgres://postgres:admin@localhost:5432/HistoricalData?sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	//defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
	return db
}
