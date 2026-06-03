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
	TotalStok   int
	TotalMasuk  int
	TotalKeluar int
}

type BarangMasuk struct {
	ID         int
	BarangID   int
	NamaBarang string
	Jumlah     int
	Tanggal    string
}

type BarangKeluar struct {
	ID         int
	BarangID   int
	NamaBarang string
	Jumlah     int
	Tanggal    string
}

type BarangMasukPageData struct {
	BarangList  []Barang
	RiwayatList []BarangMasuk
}

type BarangKeluarPageData struct {
	BarangList  []Barang
	RiwayatList []BarangKeluar
}
