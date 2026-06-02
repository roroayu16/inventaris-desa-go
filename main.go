// belajar git pertama
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Barang struct {
	ID      int
	Nama    string
	Jumlah  string
	Lokasi  string
	Kondisi string
}

type HomeData struct {
	TotalBarang int
}

var nextID = 1
var daftarBarang []Barang

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := HomeData{
		TotalBarang: len(daftarBarang),
	}

	tmpl.Execute(w, data)
}

func barangHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/barang.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, daftarBarang)
}

func tambahBarangHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nama := r.FormValue("nama")
		jumlah := r.FormValue("jumlah")
		lokasi := r.FormValue("lokasi")
		kondisi := r.FormValue("kondisi")

		barangBaru := Barang{
			ID:      nextID,
			Nama:    nama,
			Jumlah:  jumlah,
			Lokasi:  lokasi,
			Kondisi: kondisi,
		}

		nextID++
		daftarBarang = append(daftarBarang, barangBaru)

		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/tambah_barang.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func editBarangHandler(w http.ResponseWriter, r *http.Request) {

	// POST
	if r.Method == "POST" {
		idstr := r.FormValue("id")

		nama := r.FormValue("nama")
		jumlah := r.FormValue("jumlah")
		lokasi := r.FormValue("lokasi")
		kondisi := r.FormValue("kondisi")

		id, err := strconv.Atoi(idstr)

		if err != nil {
			http.Redirect(w, r, "/barang", http.StatusSeeOther)
			return
		}

		for i, barang := range daftarBarang {
			if barang.ID == id {
				daftarBarang[i].Nama = nama
				daftarBarang[i].Jumlah = jumlah
				daftarBarang[i].Lokasi = lokasi
				daftarBarang[i].Kondisi = kondisi

				break
			}
		}

		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	// GET
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	for _, barang := range daftarBarang {

		if barang.ID == id {

			tmpl, err := template.ParseFiles("templates/edit_barang.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl.Execute(w, barang)

			return
		}
	}

	fmt.Fprintf(w, "Barang tidak ditemukan")
}

func hapusBarangHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	for i, barang := range daftarBarang {
		if barang.ID == id {
			daftarBarang = append(
				daftarBarang[:i],
				daftarBarang[i+1:]...,
			)

			break
		}
	}

	http.Redirect(w, r, "/barang", http.StatusSeeOther)

}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/barang", barangHandler)
	http.HandleFunc("/barang/tambah", tambahBarangHandler)
	http.HandleFunc("/barang/edit", editBarangHandler)
	http.HandleFunc("/barang/hapus", hapusBarangHandler)

	fmt.Println("server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
