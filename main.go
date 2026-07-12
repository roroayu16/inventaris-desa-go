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

	// DASHBOARD
	http.HandleFunc("/", homeHandler)

	// BARANG
	http.HandleFunc("/barang", barangHandler)
	http.HandleFunc("/barang/tambah", tambahBarangHandler)
	http.HandleFunc("/barang/edit", editBarangHandler)
	http.HandleFunc("/barang/hapus", hapusBarangHandler)
	http.HandleFunc("/barang/detail", detailBarangHandler)

	// KATEGORI
	http.HandleFunc("/kategori", kategoriHandler)
	http.HandleFunc("/kategori/tambah", tambahKategoriHandler)
	http.HandleFunc("/kategori/edit", editKategoriHandler)
	http.HandleFunc("/kategori/hapus", hapusKategoriHandler)

	// BARANG MASUK
	http.HandleFunc("/barang-masuk", barangMasukHandler)

	// BARANG KELUAR
	http.HandleFunc("/barang-keluar", barangKeluarHandler)

	// LAPORAN
	http.HandleFunc("/export/barang", exportBarangExcelHandler)

	log.Println("server berjalan di http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
