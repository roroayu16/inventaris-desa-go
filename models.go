package main

type Barang struct {
	ID      int
	Nama    string
	Jumlah  int
	Lokasi  string
	Kondisi string
}

type HomeData struct {
	TotalBarang int
}

type BarangMasuk struct {
	ID       int
	BarangID int
	Jumlah   int
	Tanggal  string
}

type BarangKeluar struct {
	ID       int
	BarangID int
	Jumlah   int
	Tanggal  string
}
