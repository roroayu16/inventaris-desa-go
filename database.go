package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	migrateDatabase()

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

func addColumn(table, column, dataType string) {
	query := fmt.Sprintf(`
		ALTER TABLE %s
		ADD COLUMN %s %s
	`, table, column, dataType)

	_, err := db.Exec(query)

	if err == nil {
		log.Printf("Migration: %s.%s berhasil", table, column)
		return
	}

	if strings.Contains(err.Error(), "duplicate column name") {
		log.Printf("Migration: %s.%s sudah ada", table, column)
		return
	}

	log.Println(err)
}

func migrateDatabase() {

	addColumn("barang", "kategori_id", "INTEGER")
	addColumn("barang", "kode_barang", "TEXT")
	addColumn("barang_masuk", "keterangan", "TEXT")
	addColumn("barang_keluar", "diambil_oleh", "TEXT")
	addColumn("barang_keluar", "keperluan", "TEXT")
	addColumn("barang_keluar", "keterangan", "TEXT")
}

func generateKodeBarang(kategoriID string) (string, error) {

	var kodeKategori string

	err := db.QueryRow(`
		SELECT kode
		FROM kategori
		WHERE id = ?
	`, kategoriID).Scan(&kodeKategori)

	if err != nil {
		return "", err
	}

	var kodeTerakhir string

	err = db.QueryRow(`
		SELECT kode_barang
		FROM barang
		WHERE kategori_id = ?
		ORDER BY kode_barang DESC
		LIMIT 1
	`, kategoriID).Scan(&kodeTerakhir)

	if err == sql.ErrNoRows {
		return fmt.Sprintf(
			"%s-%03d",
			kodeKategori,
			1,
		), nil
	}

	if err != nil {
		return "", err
	}

	var nomor int

	fmt.Sscanf(
		kodeTerakhir,
		kodeKategori+"-%d",
		&nomor,
	)

	nomor++

	kodeBaru := fmt.Sprintf(
		"%s-%03d",
		kodeKategori,
		nomor,
	)

	return kodeBaru, nil
}

func getAllBarang() ([]Barang, error) {
	rows, err := db.Query(`
		SELECT 
			b.id,
			b.kode_barang,
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
			&b.KodeBarang,
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
		b.kode_barang,
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
		&barang.KodeBarang,
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
	kodeBarang, err := generateKodeBarang(kategoriID)

	if err != nil {
		return err
	}

	query := `
	INSERT INTO barang (
		kode_barang,
		kategori_id,
		nama,
		stok_awal,
		jumlah,
		tempat,
		kondisi
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(
		query,
		kodeBarang,
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
			b.id,
			b.kode_barang,
			k.nama,
			b.nama,
			b.tempat,
			b.kondisi,
			b.jumlah
		FROM barang b
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
		ORDER BY b.kode_barang
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
			&b.KodeBarang,
			&b.Kategori,
			&b.Nama,
			&b.Tempat,
			&b.Kondisi,
			&b.Jumlah,
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
	keterangan string,
) error {
	query := `
	INSERT INTO barang_masuk (
		barang_id,
		jumlah,
		tanggal,
		keterangan
	)
	VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		barangID,
		jumlah,
		tanggal,
		keterangan,
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
	diambilOleh string,
	keperluan string,
	keterangan string,
) error {
	query := `
	INSERT INTO barang_keluar (
		barang_id,
		jumlah,
		tanggal,
		diambil_oleh,
		keperluan,
		keterangan
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		barangID,
		jumlah,
		tanggal,
		diambilOleh,
		keperluan,
		keterangan,
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
			b.kode_barang,
			k.nama,
			b.nama,
			b.tempat,
			b.kondisi,
			bm.jumlah,
			bm.tanggal,
			bm.keterangan
		FROM barang_masuk bm
		JOIN barang b
			ON bm.barang_id = b.id
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
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
			&b.KodeBarang,
			&b.Kategori,
			&b.NamaBarang,
			&b.Tempat,
			&b.Kondisi,
			&b.Jumlah,
			&b.Tanggal,
			&b.Keterangan,
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
			bk.id,
			bk.barang_id,
			b.kode_barang,
			k.nama,
			b.nama,
			b.tempat,
			b.kondisi,
			bk.jumlah,
			bk.tanggal,
			bk.diambil_oleh,
			bk.keperluan,
			bk.keterangan
		FROM barang_keluar bk
		JOIN barang b
			ON bk.barang_id = b.id
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
		ORDER BY bk.id DESC
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
			&b.KodeBarang,
			&b.Kategori,
			&b.NamaBarang,
			&b.Tempat,
			&b.Kondisi,
			&b.Jumlah,
			&b.Tanggal,
			&b.DiambilOleh,
			&b.Keperluan,
			&b.Keterangan,
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
			bm.id,
			bm.barang_id,
			b.kode_barang,
			k.nama,
			b.nama,
			b.tempat,
			b.kondisi,
			bm.jumlah,
			bm.tanggal,
			bm.keterangan
		FROM barang_masuk bm
		JOIN barang b
			ON bm.barang_id = b.id
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
		WHERE bm.barang_id = ?
		ORDER BY bm.tanggal DESC
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
			&b.KodeBarang,
			&b.Kategori,
			&b.NamaBarang,
			&b.Tempat,
			&b.Kondisi,
			&b.Jumlah,
			&b.Tanggal,
			&b.Keterangan,
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
			bk.id,
			bk.barang_id,
			b.kode_barang,
			k.nama,
			b.nama,
			b.tempat,
			b.kondisi,
			bk.jumlah,
			bk.tanggal,
			bk.diambil_oleh,
			bk.keperluan,
			bk.keterangan
		FROM barang_keluar bk
		JOIN barang b
			ON bk.barang_id = b.id
		LEFT JOIN kategori k
			ON b.kategori_id = k.id
		WHERE bk.barang_id = ?
		ORDER BY bk.tanggal DESC
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
			&b.KodeBarang,
			&b.Kategori,
			&b.NamaBarang,
			&b.Tempat,
			&b.Kondisi,
			&b.Jumlah,
			&b.Tanggal,
			&b.DiambilOleh,
			&b.Keperluan,
			&b.Keterangan,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, b)
	}
	return data, nil
}
