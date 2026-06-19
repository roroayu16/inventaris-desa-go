package main

type Barang struct {
	ID         int
	KategoriID int
	Kategori   string

	Nama    string
	Tempat  string
	Kondisi string

	StokAwal    int
	TotalMasuk  int
	TotalKeluar int
	Jumlah      int
}

type Kategori struct {
	ID   int
	Kode string
	Nama string
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

type DetailBarang struct {
	Barang Barang

	RiwayatMasuk  []BarangMasuk
	RiwayatKeluar []BarangKeluar
}

type TambahBarangData struct {
	Kategori []Kategori
}

type EditBarangData struct {
	Barang   Barang
	Kategori []Kategori
}
