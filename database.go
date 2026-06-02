package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	log.Println("Database INIT Start")

	var err error

	db, err = sql.Open("sqlite", "./inventaris.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	log.Println("Database INIT Success")
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS barang (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nama TEXT NOT NULL,
		jumlah INTEGER NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
