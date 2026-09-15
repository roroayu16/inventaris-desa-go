package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var db *sql.DB

var ErrKategoriMasihDigunakan = errors.New("kategori masih digunakan")

// ==============================================
// DATABASE INITIALIZATION
// ==============================================
func initDB() {
	log.Println("Database INIT Start")

	var err error

	db, err = sql.Open("sqlite", "./inventaris.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
	createInitialSuperAdministrator()

	log.Println("Database INIT Success")
}

// ==============================================
// DATABASE SCHEMA
// ==============================================
func createTable() {
	// BARANG
	queryBarang := `
	CREATE TABLE IF NOT EXISTS barang (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		kode_barang TEXT NOT NULL,
		kategori_id INTEGER NOT NULL,
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

	// KATEGORI
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

	// BARANG MASUK
	queryMasuk := `
	CREATE TABLE IF NOT EXISTS barang_masuk (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		barang_id INTEGER NOT NULL,
		jumlah INTEGER NOT NULL,
		tanggal TEXT NOT NULL,
		keterangan TEXT
	);
	`
	_, err = db.Exec(queryMasuk)
	if err != nil {
		log.Fatal(err)
	}

	// BARANG KELUAR
	queryKeluar := `
	CREATE TABLE IF NOT EXISTS barang_keluar (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		barang_id INTEGER NOT NULL,
		jumlah INTEGER NOT NULL,
		tanggal TEXT NOT NULL,
		diambil_oleh TEXT,
		keperluan TEXT,
		keterangan TEXT
	);
	`
	_, err = db.Exec(queryKeluar)
	if err != nil {
		log.Fatal(err)
	}

	// USERS
	queryUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		nama TEXT NOT NULL,
		role TEXT NOT NULL CHECK (
			role IN ('administrator', 'super_administrator')
		),
		aktif INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	`

	if _, err := db.Exec(queryUsers); err != nil {
		log.Fatal(err)
	}

	// SESSIONS
	querySessions := `
	CREATE TABLE IF NOT EXISTS sessions (
		id 	INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	if _, err := db.Exec(querySessions); err != nil {
		log.Fatal(err)
	}

	// LOG KEGIATAN
	queryLogKegiatan := `
	CREATE TABLE IF NOT EXISTS log_kegiatan (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		aktivitas TEXT NOT NULL,
		keterangan TEXT,
		waktu TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	if _, err := db.Exec(queryLogKegiatan); err != nil {
		log.Fatal(err)
	}
}

// ==============================================
// USER DATABASE FUNCTIONS
// ==============================================
func getUserByUsername(username string) (User, error) {
	var user User

	err := db.QueryRow(`
		SELECT
			id,
			username,
			password,
			nama,
			role,
			aktif,
			created_at,
			updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Nama,
		&user.Role,
		&user.Aktif,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func getAllUsers() ([]User, error) {

	rows, err := db.Query(`
        SELECT
            id,
            username,
            password,
            nama,
            role,
            aktif,
            created_at,
            updated_at
        FROM users
        ORDER BY id ASC
    `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {

		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Nama,
			&user.Role,
			&user.Aktif,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ==============================================
// PASSWORD SECURITY
// ==============================================
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}

// reset Admin Password
func resetAdminPassword(newPassword string) error {

	hashedPassword, err := hashPassword(newPassword)

	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`
        UPDATE users
        SET
            password = ?,
            updated_at = ?
        WHERE username = 'admin'
    `,
		hashedPassword,
		now,
	)

	return err
}

func generateTemporaryPassword() (string, error) {

	const characters = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

	const passwordLength = 10

	randomBytes := make([]byte, passwordLength)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	password := make([]byte, passwordLength)

	for i := range randomBytes {
		password[i] = characters[int(randomBytes[i])%len(characters)]
	}

	return string(password), nil
}

// ==============================================
// SESSION MANAGEMENT
// ==============================================
func generateSessionToken() (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func createSession(userID int) (string, error) {

	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(8 * time.Hour)

	_, err = db.Exec(`
        INSERT INTO sessions (
            user_id,
            token,
            created_at,
            expires_at
        )
        VALUES (?, ?, ?, ?)
    `,
		userID,
		token,
		now.Format("2006-01-02 15:04:05"),
		expiresAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

// ==============================================
// SESSION VALIDATION
// ==============================================
func getUserBySessionToken(token string) (User, error) {

	var user User
	var expiresAt string

	err := db.QueryRow(`
        SELECT
            users.id,
            users.username,
            users.password,
            users.nama,
            users.role,
            users.aktif,
            users.created_at,
            users.updated_at,
            sessions.expires_at
        FROM sessions
        JOIN users
            ON users.id = sessions.user_id
        WHERE sessions.token = ?
    `, token).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Nama,
		&user.Role,
		&user.Aktif,
		&user.CreatedAt,
		&user.UpdatedAt,
		&expiresAt,
	)

	if err != nil {
		return User{}, err
	}

	// Cek apakah session sudah kedaluwarsa
	expirationTime, err := time.Parse(
		"2006-01-02 15:04:05",
		expiresAt,
	)

	if err != nil {
		return User{}, err
	}

	if time.Now().After(expirationTime) {
		// Hapus session yang sudah kedaluwarsa
		_, _ = db.Exec(`
            DELETE FROM sessions
            WHERE token = ?
        `, token)

		return User{}, errors.New("session expired")
	}

	// Cek apakah akun masih aktif
	if user.Aktif != 1 {
		return User{}, errors.New("user inactive")
	}

	return user, nil
}

// ==============================================
// LOGIN CHECK
// ==============================================
func getCurrentUser(r *http.Request) (User, error) {

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return User{}, err
	}

	return getUserBySessionToken(cookie.Value)
}

func getCurrentSessionInfo(r *http.Request) (SessionInfo, error) {

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return SessionInfo{}, err
	}

	var expiresAt string

	err = db.QueryRow(`
        SELECT expires_at
        FROM sessions
        WHERE token = ?
    `, cookie.Value).Scan(&expiresAt)

	if err != nil {
		return SessionInfo{}, err
	}

	expirationTime, err := time.ParseInLocation(
		"2006-01-02 15:04:05",
		expiresAt,
		time.Local,
	)

	if err != nil {
		return SessionInfo{}, err
	}

	if time.Now().After(expirationTime) {
		_, _ = db.Exec(`
            DELETE FROM sessions
            WHERE token = ?
        `, cookie.Value)

		return SessionInfo{}, errors.New("session expired")
	}

	return SessionInfo{
		ExpiresAt: expirationTime,
	}, nil
}

// ==============================================
// INITIAL ADMINISTRATOR
// ==============================================

func createInitialSuperAdministrator() {

	createUserIfNotExists(
		"superadmin",
		"superadmin123",
		"Super Administrator",
		"super_administrator",
	)

	createUserIfNotExists(
		"admin",
		"admin123",
		"Administrator",
		"administrator",
	)
}

func createUserIfNotExists(
	username string,
	password string,
	nama string,
	role string,
) {

	var count int

	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM users
        WHERE username = ?
    `, username).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		return
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`
        INSERT INTO users (
            username,
            password,
            nama,
            role,
            aktif,
            created_at,
            updated_at
        )
        VALUES (?, ?, ?, ?, 1, ?, ?)
    `,
		username,
		hashedPassword,
		nama,
		role,
		now,
		now,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("User created:", username)
}

// ==============================================
// BARANG
// ==============================================

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

	fmt.Scanf(
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

func sinkronkanKodeBarang(kategoriID string) error {
	var kodeKategori string

	err := db.QueryRow(`
		SELECT kode
		FROM kategori
		WHERE id = ?
	`, kategoriID).Scan(&kodeKategori)

	if err != nil {
		return err
	}

	rows, err := db.Query(`
		SELECT id
		FROM barang
		WHERE kategori_id = ?
		ORDER BY kode_barang
	`, kategoriID)

	if err != nil {
		return err
	}

	var barangIDs []int

	for rows.Next() {
		var barangID int

		err := rows.Scan(&barangID)

		if err != nil {
			return err
		}

		barangIDs = append(
			barangIDs,
			barangID,
		)
	}

	// Tutup rows sebelum melakukan UPDATE
	// untuk menghindari SQLITE_BUSY
	rows.Close()

	nomor := 1

	for _, barangID := range barangIDs {

		kodeBaru := fmt.Sprintf(
			"%s-%03d",
			kodeKategori,
			nomor,
		)

		_, err = db.Exec(`
			UPDATE barang
			SET kode_barang = ?
			WHERE id = ?
		`,
			kodeBaru,
			barangID,
		)

		if err != nil {
			return err
		}

		nomor++
	}

	return nil
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

// ==============================================
// KATEGORI
// ==============================================
func getAllKategori() ([]Kategori, error) {
	rows, err := db.Query(`
		SELECT id, kode, nama
		FROM kategori
		ORDER BY nama
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kategori []Kategori

	for rows.Next() {
		var k Kategori

		err := rows.Scan(
			&k.ID,
			&k.Kode,
			&k.Nama,
		)
		if err != nil {
			return nil, err
		}

		kategori = append(kategori, k)
	}

	return kategori, nil
}

func getKategoriByID(id string) (Kategori, error) {
	var kategori Kategori

	err := db.QueryRow(`
		SELECT id, kode, nama
		FROM kategori
		WHERE id = ?
	`, id).Scan(
		&kategori.ID,
		&kategori.Kode,
		&kategori.Nama,
	)

	return kategori, err
}

func insertKategori(kode, nama string) error {
	_, err := db.Exec(`
		INSERT INTO kategori (kode, nama)
		VALUES (?, ?)
	`, kode, nama)

	return err
}

func updateKategori(id, kode, nama string) error {
	_, err := db.Exec(`
		UPDATE kategori
		SET kode = ?, nama = ?
		WHERE id = ?
	`, kode, nama, id)

	if err != nil {
		return err
	}

	return sinkronkanKodeBarang(id)
}

func deleteKategori(id string) error {
	jumlah, err := getJumlahBarangByKategori(id)

	if err != nil {
		return err
	}

	if jumlah > 0 {
		return ErrKategoriMasihDigunakan
	}

	_, err = db.Exec(`
		DELETE FROM kategori
		WHERE id = ?
	`, id)

	return err
}

func getJumlahBarangByKategori(id string) (int, error) {
	var jumlah int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM barang
		WHERE kategori_id = ?
	`, id).Scan(&jumlah)

	return jumlah, err
}

// ==============================================
// DASHBOARD
// ==============================================

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

// ==============================================
// BARANG MASUK
// ==============================================

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

// ==============================================
// BARANG KELUAR
// ==============================================

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
