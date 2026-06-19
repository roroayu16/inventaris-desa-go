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
	migrateBarangKategori()

	log.Println("Database INIT Success")
}

func createTable() {
	//barang
	queryBarang := `
	CREATE TABLE IF NOT EXISTS barang (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nama TEXT NOT NULL,
		stok_awal INTEGER NOT NULL,
		jumlah INTEGER NOT NULL,
		tempat TEXT NOT NULL,
		kondisi TEXT NOT NULL
	);
	`

	_, err := db.Exec(queryBarang)
	if err != nil {
		log.Fatal(err)
	}

	//Kategori
	queryKategori := `
	CREATE TABLE IF NOT EXISTS kategori (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		kode TEXT NOT NULL UNIQUE,
		nama TEXT NOT NULL UNIQUE
	);
	`

	if _, err := db.Exec(queryKategori); err != nil {
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

func migrateBarangKategori() {

	_, err := db.Exec(`
		ALTER TABLE barang
		ADD COLUMN kategori_id INTEGER
	`)

	if err != nil {
		log.Println("Migration kategori_id dilewati:", err)
	}

	log.Println("Migration kategori_id berhasil")
}

func getAllBarang() ([]Barang, error) {
	rows, err := db.Query(`
		SELECT 
			b.id,
			k.nama AS kategori, 
			b.nama,
			b.tempat, 
			b.kondisi,
			b.stok_awal, 

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
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
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
			&b.Kategori,
			&b.Nama,
			&b.Tempat,
			&b.Kondisi,
			&b.StokAwal,
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
		b.id,
		b.kategori_id,
		k.nama,
		b.nama,
		b.stok_awal,
		b.tempat,
		b.kondisi,
		b.jumlah
	FROM barang b
	LEFT JOIN kategori k
		ON b.kategori_id = k.id
	WHERE b.id = ?
	`

	err := db.QueryRow(
		query,
		id,
	).Scan(
		&barang.ID,
		&barang.KategoriID,
		&barang.Kategori,
		&barang.Nama,
		&barang.StokAwal,
		&barang.Tempat,
		&barang.Kondisi,
		&barang.Jumlah,
	)

	return barang, err
}

func insertBarang(kategoriID, nama, jumlah, tempat, kondisi string) error {
	query := `
	INSERT INTO barang (
		kategori_id,
		nama,
		stok_awal,
		jumlah,
		tempat,
		kondisi
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		kategoriID,
		nama,
		jumlah,
		jumlah,
		tempat,
		kondisi,
	)

	return err
}

func updateBarang(
	id int,
	kategoriID string,
	nama string,
	jumlah string,
	tempat string,
	kondisi string,
) error {
	query := `
	UPDATE barang
	SET
		kategori_id = ?,
		nama = ?,
		jumlah = ?,
		tempat = ?,
		kondisi = ?
	WHERE id = ?
	`

	_, err := db.Exec(
		query,
		kategoriID,
		nama,
		jumlah,
		tempat,
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
			tempat,
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
			&b.Tempat,
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
