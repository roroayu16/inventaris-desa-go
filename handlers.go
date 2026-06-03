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

func tambahBarangHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nama := r.FormValue("nama")
		jumlah := r.FormValue("jumlah")
		lokasi := r.FormValue("lokasi")
		kondisi := r.FormValue("kondisi")

		err := insertBarang(
			nama,
			jumlah,
			lokasi,
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

	tmpl.Execute(w, nil)
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

		nama := r.FormValue("nama")
		jumlah := r.FormValue("jumlah")
		lokasi := r.FormValue("lokasi")
		kondisi := r.FormValue("kondisi")

		id, err := strconv.Atoi(idstr)

		if err != nil {
			http.Redirect(w, r, "/barang", http.StatusSeeOther)
			return
		}

		err = updateBarang(
			id,
			nama,
			jumlah,
			lokasi,
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
		http.Redirect(
			w,
			r,
			"/barang",
			http.StatusSeeOther,
		)
		return
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

	tmpl.Execute(w, barang)

	fmt.Fprintf(w, "Barang tidak ditemukan")
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

	f.SetCellValue(
		"Sheet1",
		"A1",
		"ID",
	)

	f.SetCellValue(
		"Sheet1",
		"B1",
		"Nama Barang",
	)

	f.SetCellValue(
		"Sheet1",
		"C1",
		"Jumlah",
	)

	f.SetCellValue(
		"Sheet1",
		"D1",
		"Lokasi",
	)

	f.SetCellValue(
		"Sheet1",
		"E1",
		"Kondisi",
	)

	for i, barang := range barangList {
		row := i + 2

		f.SetCellValue(
			"Sheet1",
			fmt.Sprintf("A%d", row),
			barang.ID,
		)

		f.SetCellValue(
			"Sheet1",
			fmt.Sprintf("B%d", row),
			barang.Nama,
		)

		f.SetCellValue(
			"Sheet1",
			fmt.Sprintf("C%d", row),
			barang.Jumlah,
		)

		f.SetCellValue(
			"Sheet1",
			fmt.Sprintf("D%d", row),
			barang.Lokasi,
		)

		f.SetCellValue(
			"Sheet1",
			fmt.Sprintf("E%d", row),
			barang.Kondisi,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="inventaris_barang.xlsx`,
	)

	err = f.Write(w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
