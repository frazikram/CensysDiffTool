package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_ip TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		data TEXT NOT NULL,
		UNIQUE(host_ip, timestamp)
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	fmt.Println("Database initialized:", filepath)
	return db
}
