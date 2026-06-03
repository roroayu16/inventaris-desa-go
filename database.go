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
		jumlah TEXT NOT NULL,
		lokasi TEXT NOT NULL,
		kondisi TEXT NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllBarang() ([]Barang, error) {
	rows, err := db.Query(`
		SELECT id, nama, jumlah, lokasi, kondisi
		FROM barang
		ORDER BY id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var barangList []Barang

	for rows.Next() {
		var b Barang
		err := rows.Scan(
			&b.ID,
			&b.Nama,
			&b.Jumlah,
			&b.Lokasi,
			&b.Kondisi,
		)

		if err != nil {
			return nil, err
		}

		barangList = append(barangList, b)
	}

	return barangList, nil
}

func insertBarang(nama, jumlah, lokasi, kondisi string) error {
	query := `
	INSERT INTO barang (
		nama,
		jumlah,
		lokasi,
		kondisi
	)
	VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		nama,
		jumlah,
		lokasi,
		kondisi,
	)

	return err
}

func getBarangByID(id int) (Barang, error) {
	var barang Barang

	query := `
	SELECT
		id,
		nama,
		jumlah,
		lokasi,
		kondisi
	FROM barang
	WHERE id = ?
	`

	err := db.QueryRow(
		query,
		id,
	).Scan(
		&barang.ID,
		&barang.Nama,
		&barang.Jumlah,
		&barang.Lokasi,
		&barang.Kondisi,
	)

	return barang, err
}

func updateBarang(
	id int,
	nama string,
	jumlah string,
	lokasi string,
	kondisi string,
) error {
	query := `
	UPDATE barang
	SET
		nama = ?,
		jumlah = ?,
		lokasi = ?,
		kondisi = ?
	WHERE id = ?
	`

	_, err := db.Exec(
		query,
		nama,
		jumlah,
		lokasi,
		kondisi,
		id,
	)

	return err
}

func deleteBarang(id int) error {
	query := `
	DELETE FROM barang
	WHERE id = ?
	`

	_, err := db.Exec(
		query,
		id,
	)

	return err
}

func getTotalBarang() (int, error) {
	var total int

	query := `
	SELECT COUNT(*)
	FROM barang
	`

	err := db.QueryRow(query).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}
