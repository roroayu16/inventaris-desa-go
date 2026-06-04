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
	//barang
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

	//barang-masuk
	queryMasuk := `
	CREATE TABLE IF NOT EXISTS barang_masuk (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		barang_id INTEGER NOT NULL,
		jumlah INTEGER NOT NULL,
		tanggal TEXT NOT NULL
	);
	`
	_, err = db.Exec(queryMasuk)

	if err != nil {
		log.Fatal(err)
	}

	//barang-keluar
	queryKeluar := `
	CREATE TABLE IF NOT EXISTS barang_keluar (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		barang_id INTEGER NOT NULL,
		jumlah INTEGER NOT NULL,
		tanggal TEXT NOT NULL
	);
	`
	_, err = db.Exec(queryKeluar)

	if err != nil {
		log.Fatal(err)
	}
}

func getAllBarang() ([]Barang, error) {
	rows, err := db.Query(`
		SELECT 
			b.id, 
			b.nama, 
			b.lokasi, 
			b.kondisi,

			COALESCE(
				(
					SELECT SUM(jumlah)
					FROM barang_masuk bm
					WHERE bm.barang_id = b.id
				),
				0
			) AS total_masuk,
			
			COALESCE(
				(
					SELECT SUM(jumlah)
					FROM barang_keluar bk
					WHERE bk.barang_id = b.id
				),
				0
			) AS total_keluar,
			
			b.jumlah

		FROM barang b
		ORDER BY b.id
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
			&b.Lokasi,
			&b.Kondisi,
			&b.TotalMasuk,
			&b.TotalKeluar,
			&b.Jumlah,
		)

		if err != nil {
			return nil, err
		}

		barangList = append(barangList, b)
	}

	return barangList, nil
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

func getTotalStok() (int, error) {
	var total int

	query := `
	SELECT COALESCE(SUM(jumlah), 0)
	FROM barang
	`

	err := db.QueryRow(query).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func getTotalBarangMasuk() (int, error) {
	var total int

	query := `
	SELECT COALESCE(SUM(jumlah), 0)
	FROM barang_masuk
	`

	err := db.QueryRow(query).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func getTotalBarangKeluar() (int, error) {
	var total int

	query := `
	SELECT COALESCE(SUM(jumlah), 0)
	FROM barang_keluar
	`

	err := db.QueryRow(query).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func getAllBarangForDropDown() ([]Barang, error) {
	rows, err := db.Query(`
		SELECT
			id,
			nama,
			jumlah,
			lokasi,
			kondisi
		FROM barang
		ORDER BY nama
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

func insertBarangMasuk(
	barangID int,
	jumlah int,
	tanggal string,
) error {
	query := `
	INSERT INTO barang_masuk (
		barang_id,
		jumlah,
		tanggal
	)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(
		query,
		barangID,
		jumlah,
		tanggal,
	)

	return err
}

func updateStokMasuk(
	barangID int,
	jumlah int,
) error {
	query := `
	UPDATE barang
	SET jumlah = jumlah + ?
	WHERE id = ?
	`

	_, err := db.Exec(
		query,
		jumlah,
		barangID,
	)

	return err
}

func insertBarangKeluar(
	barangID int,
	jumlah int,
	tanggal string,
) error {
	query := `
	INSERT INTO barang_keluar (
		barang_id,
		jumlah,
		tanggal
	)
	VALUES (?, ?, ?)
	`

	_, err := db.Exec(
		query,
		barangID,
		jumlah,
		tanggal,
	)

	return err
}

func updateStokKeluar(
	barangID int,
	jumlah int,
) error {
	query := `
	UPDATE barang
	SET jumlah = jumlah - ?
	WHERE id = ?
	`

	_, err := db.Exec(
		query,
		jumlah,
		barangID,
	)

	return err
}

func getAllBarangMasuk() ([]BarangMasuk, error) {
	rows, err := db.Query(`
		SELECT
			bm.id,
			bm.barang_id,
			b.nama,
			bm.jumlah,
			bm.tanggal
		FROM barang_masuk bm
		JOIN barang b
		ON bm.barang_id = b.id
		ORDER BY bm.id DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []BarangMasuk

	for rows.Next() {
		var b BarangMasuk
		err := rows.Scan(
			&b.ID,
			&b.BarangID,
			&b.NamaBarang,
			&b.Jumlah,
			&b.Tanggal,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, b)
	}

	return data, nil
}

func getAllBarangKeluar() ([]BarangKeluar, error) {
	rows, err := db.Query(`
		SELECT
			bm.id,
			bm.barang_id,
			b.nama,
			bm.jumlah,
			bm.tanggal
		FROM barang_keluar bm
		JOIN barang b
		ON bm.barang_id = b.id
		ORDER BY bm.id DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []BarangKeluar

	for rows.Next() {
		var b BarangKeluar
		err := rows.Scan(
			&b.ID,
			&b.BarangID,
			&b.NamaBarang,
			&b.Jumlah,
			&b.Tanggal,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, b)
	}

	return data, nil
}

func getBarangMasukByBarangID(barangID int) ([]BarangMasuk, error) {
	rows, err := db.Query(`
		SELECT
			id,
			barang_id,
			jumlah,
			tanggal
		FROM barang_masuk
		WHERE barang_id = ?
		ORDER BY tanggal DESC
	`, barangID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var data []BarangMasuk

	for rows.Next() {
		var b BarangMasuk

		err := rows.Scan(
			&b.ID,
			&b.BarangID,
			&b.Jumlah,
			&b.Tanggal,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, b)
	}
	return data, nil
}

func getBarangKeluarByBarangID(barangID int) ([]BarangKeluar, error) {
	rows, err := db.Query(`
		SELECT
			id,
			barang_id,
			jumlah,
			tanggal
		FROM barang_keluar
		WHERE barang_id = ?
		ORDER BY tanggal DESC
	`, barangID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var data []BarangKeluar

	for rows.Next() {
		var b BarangKeluar

		err := rows.Scan(
			&b.ID,
			&b.BarangID,
			&b.Jumlah,
			&b.Tanggal,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, b)
	}
	return data, nil
}
