# SIPACAR

**Sistem Informasi Pengelolaan Inventaris Carangrejo**

SIPACAR merupakan aplikasi berbasis web yang dikembangkan untuk membantu Pemerintah Desa Carangrejo dalam mengelola inventaris operasional desa secara lebih tertata, terdokumentasi, dan mudah dipantau.

Aplikasi ini digunakan untuk mencatat berbagai barang inventaris seperti Alat Tulis Kantor (ATK), laptop, printer, meja, kursi, lemari, serta peralatan operasional lainnya yang digunakan dalam aktivitas pemerintahan desa.

> **Catatan**
>
> SIPACAR **bukan pengganti SIPADES**. SIPACAR berfungsi sebagai sistem pendukung yang berfokus pada pengelolaan inventaris internal desa, sedangkan SIPADES tetap digunakan sebagai sistem administrasi aset pemerintah desa.

---

# Philosophy

SIPACAR dikembangkan dengan beberapa prinsip utama:

* Sederhana dan mudah digunakan.
* Ringan serta tidak memerlukan server database terpisah.
* Mudah dipelihara dan dikembangkan.
* Mengutamakan keterbacaan kode (readability).
* Menggunakan komponen yang reusable.
* Memisahkan business logic dan presentation layer.
* Cocok digunakan sebagai aplikasi inventaris internal desa.

---

# Features

## Dashboard

* Menampilkan Total Barang
* Menampilkan Total Stok
* Menampilkan Total Barang Masuk
* Menampilkan Total Barang Keluar

---

## Data Barang

* CRUD Barang
* Kode barang otomatis berdasarkan kategori
* Informasi stok awal
* Informasi stok akhir
* Detail barang
* Riwayat transaksi barang

---

## Kategori

* CRUD Kategori
* Kode kategori unik
* Sinkronisasi otomatis kode barang ketika kode kategori berubah
* Validasi agar kategori yang masih digunakan tidak dapat dihapus

---

## Barang Masuk

* Pencatatan transaksi barang masuk
* Riwayat transaksi
* Penambahan stok otomatis

---

## Barang Keluar

* Pencatatan transaksi barang keluar
* Riwayat transaksi
* Pengurangan stok otomatis
* Validasi stok agar tidak menjadi negatif

---

## Laporan

* Export Data Inventaris ke Microsoft Excel (.xlsx)

---

## User Interface

* Responsive Navbar
* Bootstrap Modal
* Flash Message
* Bootstrap 5
* CSS Global
* Reusable Template

---

# Technology Stack

* Go (Golang)
* SQLite
* Bootstrap 5.3.7
* HTML Template
* Excelize

---

# Project Structure

```text
SIPACAR
│
├── main.go
├── handlers.go
├── database.go
├── models.go
├── template.go
├── flash.go
│
├── templates
│   ├── home.html
│   ├── barang.html
│   ├── detail_barang.html
│   ├── tambah_barang.html
│   ├── edit_barang.html
│   ├── kategori.html
│   ├── tambah_kategori.html
│   ├── edit_kategori.html
│   ├── barang_masuk.html
│   ├── barang_keluar.html
│   └── partials
│       ├── navbar.html
│       ├── footer.html
│       └── flash.html
│
├── static
│   └── css
│       └── style.css
│
└── inventaris.db
```

---

# Database Schema

SIPACAR menggunakan empat tabel utama.

## barang

| Field       |
| ----------- |
| id          |
| kode_barang |
| kategori_id |
| nama        |
| stok_awal   |
| jumlah      |
| tempat      |
| kondisi     |

---

## kategori

| Field |
| ----- |
| id    |
| kode  |
| nama  |

---

## barang_masuk

| Field      |
| ---------- |
| id         |
| barang_id  |
| jumlah     |
| tanggal    |
| keterangan |

---

## barang_keluar

| Field        |
| ------------ |
| id           |
| barang_id    |
| jumlah       |
| tanggal      |
| diambil_oleh |
| keperluan    |
| keterangan   |

---

# Application Architecture

```text
Browser
    │
    ▼
HTTP Handler
    │
    ▼
Business Logic (database.go)
    │
    ▼
SQLite Database
    │
    ▼
HTML Template
```

---

# How to Run

Clone repository

```bash
git clone <repository-url>
```

Masuk ke folder project

```bash
cd SIPACAR
```

Install dependency

```bash
go mod tidy
```

Jalankan aplikasi

```bash
go run .
```

Buka browser

```
http://localhost:8080
```

---

# Design System

## Background

```
#eef6ef
```

## Framework

* Bootstrap 5.3.7

## Notification

* Flash Message berbasis Cookie
* Auto Close
* Manual Close

---

# Current Version

**Version 1.0**

Fitur utama yang telah selesai:

* Dashboard
* CRUD Barang
* CRUD Kategori
* Barang Masuk
* Barang Keluar
* Export Excel
* Flash Message
* Responsive Navbar
* Clean Code Refactoring
* Final Database Schema

---

# Roadmap

## Version 1.1

Rencana pengembangan selanjutnya:

* Login
* Hak Akses (Role)
* Dashboard Statistik
* Export PDF
* Satuan Barang
* Stok Minimum
* Notifikasi
* Profil SIPACAR

---

# Development Notes

Beberapa prinsip yang diterapkan selama pengembangan:

* Business Logic dipisahkan dari HTTP Handler.
* Query database ditempatkan pada layer database.
* Template menggunakan reusable partial.
* CSS menggunakan global stylesheet.
* Database schema telah difinalisasi tanpa migration sementara.
* Kode dibuat dengan mempertimbangkan maintainability dan readability.

---

# Author
**Fasha**

Dikembangkan sebagai proyek Sistem Informasi Pengelolaan Inventaris Desa Carangrejo menggunakan Go (Golang), SQLite, dan Bootstrap.

---

# License

Project ini dikembangkan untuk tujuan akademik dan implementasi sistem inventaris internal Pemerintah Desa Carangrejo.