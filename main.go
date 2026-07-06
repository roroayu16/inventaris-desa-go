package main

import (
	"fmt"
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

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/barang", barangHandler)

	http.HandleFunc("/kategori", kategoriHandler)
	http.HandleFunc("/kategori/tambah", tambahKategoriHandler)
	http.HandleFunc("/kategori/edit", editKategoriHandler)
	http.HandleFunc("/kategori/hapus", hapusKategoriHandler)

	http.HandleFunc("/barang/detail", detailBarangHandler)
	http.HandleFunc("/barang/tambah", tambahBarangHandler)
	http.HandleFunc("/export/barang", exportBarangExcelHandler)

	http.HandleFunc("/barang-masuk", barangMasukHandler)
	http.HandleFunc("/barang-keluar", barangKeluarHandler)

	http.HandleFunc("/barang/edit", editBarangHandler)
	http.HandleFunc("/barang/hapus", hapusBarangHandler)

	fmt.Println("server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
