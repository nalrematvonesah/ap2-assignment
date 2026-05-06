package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgres() *sql.DB {
	connStr := "host=postgres port=5432 user=postgres password=postgres dbname=payments sslmode=disable"

	var db *sql.DB
	var err error

	// retry connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		log.Println("Retrying PostgreSQL connection...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect postgres:", err)
	}

	log.Println("Connected to PostgreSQL")

	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	paymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		order_id INT,
		amount FLOAT,
		email TEXT,
		status TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	)
	`

	_, err := db.Exec(paymentsTable)
	if err != nil {
		log.Fatal(err)
	}
}
