# MajooTest-Golang

Deskripsi singkat
- Projekt ini adalah implementasi backend sederhana menggunakan Go (Golang) untuk keperluan Skill Test. README ini menjelaskan cara menyiapkan, membangun, menjalankan, dan kontribusi dasar.

Fitur
- Struktur proyek Go modular
- Contoh endpoint / handler (HTTP)
- Pengujian unit dasar
- Dockerfile untuk containerisasi (jika tersedia)

Persyaratan
- Go 1.18+ (disarankan memakai versi terbaru yang kompatibel)
- Git
- (Opsional) Docker dan Docker Compose

Instalasi & Persiapan
1. Clone repository:
    git clone (https://github.com/edwinjordan/MajooTest-Golang.git)
2. Masuk ke direktori proyek:
    cd MajooTest-Golang
3. Unduh dependensi:
    go mod tidy

Membangun dan Menjalankan
- Menjalankan aplikasi secara lokal:
  go run main.go

- Menggunakan Docker
```bash
wsl -d Ubuntu

docker build -t majoo .
```  

- Menjalankan dengan Docker (jika tersedia Dockerfile):
```bash
docker run --rm -p 8000:8000 \
  -e APP_HOST=0.0.0.0 \
  -e APP_PORT=8000 \
  -e DATABASE_URL="postgres://postgres:aero1996@host.docker.internal:5432/zogtest-golang" \
  zogtest
```

Konfigurasi
- Variabel lingkungan (contoh):
  - PORT: port HTTP (default 8080)
  - DATABASE_URL: string koneksi database (jika digunakan)
- Cara mengatur:
  export PORT=8080
  export DATABASE_URL="postgres://user:pass@localhost:5432/dbname"

Pengujian
- Menjalankan unit test:
  go test ./... -v

Struktur Proyek (contoh)
- cmd/          - entri aplikasi (main)
- internal/     - kode aplikasi yang tidak diekspor
- database/     - konfigurasi database
- configs/      - konfigurasi
- domain/       - konfigurasi entity model tabel
- migrations/   - migration database yang akan di store ke database
- seeds/        - untuk store data ke database
- service       - untuk proses transaksi    
- Dockerfile    - definisi image jika ada
- go.mod, go.sum

Contributing
- Buat branch baru: git checkout -b feat/nama-fitur
- Buat commit kecil dan jelas
- Ajukan pull request dengan deskripsi perubahan
- Sertakan test untuk fitur/fix yang relevan

Catatan Pengembangan
- Gunakan context untuk manajemen request/timeout
- Tangani error secara eksplisit, jangan panic di handler
- Tulis test unit untuk logika bisnis utama

License
- Tambahkan file LICENSE sesuai lisensi yang dipilih (mis. MIT) jika perlu.

Kontak
- Sertakan informasi kontak atau referensi issue tracker pada repo untuk pertanyaan dan pelaporan bug.

Catatan akhir
- Sesuaikan nama package, path cmd, dan variabel lingkungan sesuai implementasi aktual di proyek.
- Dokumentasi ini bersifat generik — tambahkan detail endpoint, skema database, dan contoh request/response sesuai implementasi nyata.