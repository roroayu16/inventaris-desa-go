package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	totalBarang, err := getTotalBarang()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalStok, err := getTotalStok()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalMasuk, err := getTotalBarangMasuk()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalKeluar, err := getTotalBarangKeluar()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := HomeData{
		TotalBarang: totalBarang,
		TotalStok:   totalStok,
		TotalMasuk:  totalMasuk,
		TotalKeluar: totalKeluar,
	}

	tmpl, err := template.ParseFiles(
		"templates/home.html",
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl.Execute(w, data)
}

func barangHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/barang.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	barangList, err := getAllBarang()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, barangList)
}

func kategoriHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, kode, nama
		FROM kategori
		ORDER BY nama
	`)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
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
			http.Error(w, err.Error(), 500)
			return
		}

		kategori = append(kategori, k)
	}

	t, err := template.ParseFiles(
		"templates/kategori.html",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, kategori)
}

func tambahKategoriHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		kode := r.FormValue("kode")
		nama := r.FormValue("nama")

		_, err := db.Exec(`
			INSERT INTO kategori (kode, nama)
			VALUES (?, ?)
		`, kode, nama)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(
			w,
			r,
			"/kategori",
			http.StatusSeeOther,
		)

		return
	}

	t, err := template.ParseFiles(
		"templates/tambah_kategori.html",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, nil)
}

func editKategoriHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if r.Method == "POST" {

		kode := r.FormValue("kode")
		nama := r.FormValue("nama")

		_, err := db.Exec(`
			UPDATE kategori
			SET kode = ?, nama = ?
			WHERE id = ?
		`,
			kode,
			nama,
			id,
		)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(
			w,
			r,
			"/kategori",
			http.StatusSeeOther,
		)

		return
	}

	var kategori Kategori

	err := db.QueryRow(`
		SELECT id, kode, nama
		FROM kategori
		WHERE id = ?
	`,
		id,
	).Scan(
		&kategori.ID,
		&kategori.Kode,
		&kategori.Nama,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles(
		"templates/edit_kategori.html",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, kategori)
}

func hapusKategoriHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	_, err := db.Exec(`
		DELETE FROM kategori
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/kategori",
		http.StatusSeeOther,
	)
}

func tambahBarangHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		kategoriID := r.FormValue("kategori_id")
		nama := r.FormValue("nama")
		jumlah := r.FormValue("jumlah")
		tempat := r.FormValue("tempat")
		kondisi := r.FormValue("kondisi")

		err := insertBarang(
			kategoriID,
			nama,
			jumlah,
			tempat,
			kondisi,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/tambah_barang.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, kode, nama
		FROM kategori
		ORDER BY nama
	`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		kategori = append(kategori, k)
	}

	data := TambahBarangData{
		Kategori: kategori,
	}

	tmpl.Execute(w, data)
}

func barangMasukHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	//POST
	if r.Method == "POST" {
		barangIDStr := r.FormValue("barang_id")
		jumlahStr := r.FormValue("jumlah")
		tanggal := r.FormValue("tanggal")

		barangID, err := strconv.Atoi(barangIDStr)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		jumlah, err := strconv.Atoi(jumlahStr)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		err = insertBarangMasuk(
			barangID,
			jumlah,
			tanggal,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		err = updateStokMasuk(
			barangID,
			jumlah,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(
			w, r, "/barang", http.StatusSeeOther,
		)
		return
	}

	//GET
	barangList, err := getAllBarangForDropDown()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	riwayatList, err := getAllBarangMasuk()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := BarangMasukPageData{
		BarangList:  barangList,
		RiwayatList: riwayatList,
	}

	tmpl, err := template.ParseFiles("templates/barang_masuk.html")

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl.Execute(w, data)
}

func barangKeluarHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	//POST
	if r.Method == "POST" {
		barangIDStr := r.FormValue("barang_id")
		jumlahStr := r.FormValue("jumlah")
		tanggal := r.FormValue("tanggal")

		barangID, err := strconv.Atoi(barangIDStr)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		jumlah, err := strconv.Atoi(jumlahStr)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		barang, err := getBarangByID(barangID)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		if jumlah > barang.Jumlah {
			http.Error(
				w,
				"Jumlah keluar melebihi stok tersedia",
				http.StatusBadRequest,
			)
			return
		}

		err = insertBarangKeluar(
			barangID,
			jumlah,
			tanggal,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		err = updateStokKeluar(
			barangID,
			jumlah,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(
			w, r, "/barang", http.StatusSeeOther,
		)
		return
	}

	//GET
	barangList, err := getAllBarangForDropDown()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	riwayatList, err := getAllBarangKeluar()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := BarangKeluarPageData{
		BarangList:  barangList,
		RiwayatList: riwayatList,
	}

	tmpl, err := template.ParseFiles("templates/barang_keluar.html")

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl.Execute(w, data)
}

func editBarangHandler(w http.ResponseWriter, r *http.Request) {

	// POST
	if r.Method == "POST" {
		idstr := r.FormValue("id")

		KategoriID := r.FormValue("kategori_id")
		nama := r.FormValue("nama")
		tempat := r.FormValue("tempat")
		kondisi := r.FormValue("kondisi")

		id, err := strconv.Atoi(idstr)

		if err != nil {
			http.Redirect(w, r, "/barang", http.StatusSeeOther)
			return
		}

		barangLama, err := getBarangByID(id)

		if err != nil {
			http.Redirect(
				w, r, "/barang", http.StatusSeeOther,
			)
			return
		}

		err = updateBarang(
			id,
			KategoriID,
			nama,
			strconv.Itoa(barangLama.Jumlah),
			tempat,
			kondisi,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
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

	barang, err := getBarangByID(id)

	if err != nil {
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	rows, err := db.Query(`
		SELECT id, kode, nama
		FROM kategori
		ORDER BY nama
	`)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	var kategoriList []Kategori

	for rows.Next() {

		var k Kategori

		err := rows.Scan(
			&k.ID,
			&k.Kode,
			&k.Nama,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		kategoriList = append(
			kategoriList,
			k,
		)
	}

	data := EditBarangData{
		Barang:   barang,
		Kategori: kategoriList,
	}

	tmpl, err := template.ParseFiles(
		"templates/edit_barang.html",
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl.Execute(w, data)

}

func hapusBarangHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	err = deleteBarang(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(w, r, "/barang", http.StatusSeeOther)

}

func exportBarangExcelHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	barangList, err := getAllBarang()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f := excelize.NewFile()

	sheetName := "Data Inventaris SIPACAR"
	f.SetSheetName("Sheet1", sheetName)
	f.SetCellValue(sheetName, "A1", "Kode Barang")
	f.SetCellValue(sheetName, "B1", "Kategori")
	f.SetCellValue(sheetName, "C1", "Nama Barang")
	f.SetCellValue(sheetName, "D1", "Tempat")
	f.SetCellValue(sheetName, "E1", "Kondisi")
	f.SetCellValue(sheetName, "F1", "Stok Awal")
	f.SetCellValue(sheetName, "G1", "Stok Akhir")

	for i, barang := range barangList {
		row := i + 2

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("A%d", row),
			barang.KodeBarang,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("B%d", row),
			barang.Kategori,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("C%d", row),
			barang.Nama,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("D%d", row),
			barang.Tempat,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("E%d", row),
			barang.Kondisi,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("F%d", row),
			barang.StokAwal,
		)

		f.SetCellValue(
			sheetName,
			fmt.Sprintf("G%d", row),
			barang.Jumlah,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=inventaris_barang.xlsx",
	)

	err = f.Write(w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func detailBarangHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(
			w,
			"ID tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	barang, err := getBarangByID(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	riwayatMasuk, err := getBarangMasukByBarangID(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	riwayatKeluar, err := getBarangKeluarByBarangID(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := DetailBarang{
		Barang:        barang,
		RiwayatMasuk:  riwayatMasuk,
		RiwayatKeluar: riwayatKeluar,
	}

	tmpl := template.Must(
		template.ParseFiles(
			"templates/detail_barang.html",
		),
	)

	tmpl.Execute(w, data)
}
