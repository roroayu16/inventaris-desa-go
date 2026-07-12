package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ==============================================
// DASHBOARD
// ==============================================

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

	renderTemplate(w, r, "home.html", "dashboard", data)
}

// ==============================================
// BARANG
// ==============================================

func barangHandler(w http.ResponseWriter, r *http.Request) {
	barangList, err := getAllBarang()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "barang.html", "barang", barangList)
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

		SetFlash(w, "success", "Barang berhasil ditambahkan")

		http.Redirect(w, r, "/barang", http.StatusSeeOther)
		return
	}

	kategori, err := getAllKategori()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TambahBarangData{
		Kategori: kategori,
	}

	renderTemplate(w, r, "tambah_barang.html", "", data)
}

func editBarangHandler(w http.ResponseWriter, r *http.Request) {

	// POST
	if r.Method == "POST" {
		idstr := r.FormValue("id")

		kategoriID := r.FormValue("kategori_id")
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
			kategoriID,
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

		SetFlash(w, "success", "Barang berhasil diperbarui")

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

	kategoriList, err := getAllKategori()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := EditBarangData{
		Barang:   barang,
		Kategori: kategoriList,
	}

	renderTemplate(w, r, "edit_barang.html", "", data)

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

	SetFlash(w, "success", "Barang berhasil dihapus")

	http.Redirect(w, r, "/barang", http.StatusSeeOther)

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

	renderTemplate(w, r, "detail_barang.html", "barang", data)
}

// ==============================================
// KATEGORI
// ==============================================

func kategoriHandler(w http.ResponseWriter, r *http.Request) {
	kategori, err := getAllKategori()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "kategori.html", "kategori", kategori)
}

func tambahKategoriHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		kode := r.FormValue("kode")
		nama := r.FormValue("nama")

		err := insertKategori(kode, nama)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetFlash(w, "success", "Kategori berhasil ditambahkan")

		http.Redirect(w, r, "/kategori", http.StatusSeeOther)

		return
	}

	renderTemplate(w, r, "tambah_kategori.html", "", nil)
}

func editKategoriHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if r.Method == "POST" {

		kode := r.FormValue("kode")
		nama := r.FormValue("nama")

		err := updateKategori(id, kode, nama)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetFlash(w, "success", "Kategori berhasil diperbarui")

		http.Redirect(w, r, "/kategori", http.StatusSeeOther)

		return
	}

	var kategori Kategori

	kategori, err := getKategoriByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "edit_kategori.html", "", kategori)
}

func hapusKategoriHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	err := deleteKategori(id)

	if err != nil {
		if errors.Is(err, ErrKategoriMasihDigunakan) {
			SetFlash(w, "warning",
				"Kategori tidak dapat dihapus karena masih digunakan oleh satu atau lebih barang",
			)
			http.Redirect(w, r, "/kategori", http.StatusSeeOther)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetFlash(w, "success", "Kategori berhasil dihapus")

	http.Redirect(w, r, "/kategori", http.StatusSeeOther)
}

// ==============================================
// BARANG MASUK
// ==============================================

func barangMasukHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	//POST
	if r.Method == "POST" {
		barangIDStr := r.FormValue("barang_id")
		jumlahStr := r.FormValue("jumlah")
		tanggal := r.FormValue("tanggal")
		keterangan := r.FormValue("keterangan")

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
			keterangan,
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

		SetFlash(w, "success", "Data barang masuk berhasil dicatat")
		http.Redirect(w, r, "/barang-masuk", http.StatusSeeOther)
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

	renderTemplate(w, r, "barang_masuk.html", "barang-masuk", data)
}

// ==============================================
// BARANG KELUAR
// ==============================================

func barangKeluarHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	//POST
	if r.Method == "POST" {
		barangIDStr := r.FormValue("barang_id")
		jumlahStr := r.FormValue("jumlah")
		tanggal := r.FormValue("tanggal")
		diambilOleh := r.FormValue("diambil_oleh")
		keperluan := r.FormValue("keperluan")
		keterangan := r.FormValue("keterangan")

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

			SetFlash(w, "warning", "Jumlah barang keluar tidak boleh melebihi stok yang tersedia")
			http.Redirect(w, r, "/barang-keluar", http.StatusSeeOther)
			return
		}

		err = insertBarangKeluar(
			barangID,
			jumlah,
			tanggal,
			diambilOleh,
			keperluan,
			keterangan,
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

		SetFlash(w, "success", "Data barang keluar berhasil dicatat")
		http.Redirect(w, r, "/barang-keluar", http.StatusSeeOther)
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

	renderTemplate(w, r, "barang_keluar.html", "barang-keluar", data)
}

// ==============================================
// LAPORAN
// ==============================================

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
