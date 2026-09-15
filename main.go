package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()

	http.Handle(
		"/static/", http.StripPrefix(
			"/static/", http.FileServer(
				http.Dir("static"),
			),
		),
	)

	// LOGIN
	http.HandleFunc("/login", loginHandler)

	// DASHBOARD
	http.HandleFunc("/", homeHandler)

	// PROFIL
	http.HandleFunc("/profil", requireLogin(profilHandler))

	// KELOLA USER
	http.HandleFunc("/kelola-user", requireSuperAdmin(kelolaUserHandler))

	// RESET PASSWORD ADMIN
	http.HandleFunc("/kelola-user/reset-password", requireSuperAdmin(resetAdminPasswordHandler))

	// LOGOUT
	http.HandleFunc("/logout", logoutHandler)

	// BARANG
	http.HandleFunc("/barang", requireLogin(barangHandler))
	http.HandleFunc("/barang/tambah", requireLogin(tambahBarangHandler))
	http.HandleFunc("/barang/edit", requireLogin(editBarangHandler))
	http.HandleFunc("/barang/hapus", requireLogin(hapusBarangHandler))
	http.HandleFunc("/barang/detail", requireLogin(detailBarangHandler))

	// KATEGORI
	http.HandleFunc("/kategori", requireLogin(kategoriHandler))
	http.HandleFunc("/kategori/tambah", requireLogin(tambahKategoriHandler))
	http.HandleFunc("/kategori/edit", requireLogin(editKategoriHandler))
	http.HandleFunc("/kategori/hapus", requireLogin(hapusKategoriHandler))

	// BARANG MASUK
	http.HandleFunc("/barang-masuk", requireLogin(barangMasukHandler))

	// BARANG KELUAR
	http.HandleFunc("/barang-keluar", requireLogin(barangKeluarHandler))

	// LAPORAN
	http.HandleFunc("/export/barang", requireLogin(exportBarangExcelHandler))

	log.Println("server berjalan di http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
