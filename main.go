package main

import (
	"fmt"
	"net/http"
)

func main() {
	initDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/barang", barangHandler)
	http.HandleFunc("/barang/tambah", tambahBarangHandler)
	http.HandleFunc("/barang/edit", editBarangHandler)
	http.HandleFunc("/barang/hapus", hapusBarangHandler)

	fmt.Println("server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
