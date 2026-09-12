# Kane ga doko.

Aplikasi pencatatan piutang dan hutang dengan desain antarmuka bergaya Apple HIG (iOS).
Proyek ini dibangun menggunakan arsitektur monorepo yang berisi:
- **Backend:** Golang + SQLite (menggunakan `net/http`).
- **Frontend:** SvelteKit 5 (Vite) + Tailwind CSS + Lucide Icons.

## Fitur Utama
- **Dashboard:** Ringkasan total piutang, hutang, dan riwayat yang sudah lunas.
- **Pencatatan Utang/Piutang:** Antarmuka intuitif bergaya iOS untuk mencatat pinjaman.
- **Daftar Kontak:** Sistem pertemanan sepihak untuk menyaring daftar pilihan pengguna saat mencatat pinjaman.
- **Persetujuan Pembayaran:** Konfirmasi atau tolak pelunasan. Dilengkapi dengan efek suara khas Apple Pay saat pembayaran diterima.
- **Manajemen Pengguna (Admin):** Sistem pendaftaran mandiri yang membutuhkan persetujuan Admin sebelum pengguna baru dapat masuk (login). Admin juga dapat mengubah *password* atau menghapus pengguna.
- **Deployment-Ready:** Tersedia `Dockerfile.backend` dan `Dockerfile.frontend` yang sudah dikonfigurasi khusus untuk *environment* produksi seperti Dokploy.

## Cara Menjalankan di Lokal (Development)

### 1. Menjalankan Backend (Golang)
Buka terminal baru:
```bash
cd backend
go mod download
go run main.go
```
*Backend akan berjalan di port `8080`. Saat pertama kali dijalankan, sistem akan membuat file `database.sqlite` secara otomatis dan mendaftarkan akun admin default (Username: `admin`, Password: `admin123`).*

### 2. Menjalankan Frontend (SvelteKit)
Buka terminal baru:
```bash
cd frontend
npm install
npm run dev
```
*Frontend akan berjalan secara default di port `5173`. Silakan akses `http://localhost:5173` di browser Anda.*

## Panduan Deployment (Dokploy)
Proyek ini sangat mudah di-deploy karena sudah menyertakan file konfigurasi Docker.

1. **Deploy Backend**: 
   - Di Dokploy, buat aplikasi baru.
   - Pilih Build Type **Dockerfile** dan isi File Path dengan `Dockerfile.backend`.
   - Buka tab Ports dan ekspos port `8080`.
   - **Penting:** Buka tab Volumes, lalu buat Host Volume ke container `/app` agar data file `database.sqlite` tetap aman jika container di-restart.

2. **Deploy Frontend**: 
   - Buat aplikasi baru.
   - Pilih Build Type **Dockerfile** dan isi File Path dengan `Dockerfile.frontend`.
   - Buka tab Ports dan ekspos port `3000`.
   - Jika Anda menggunakan domain kustom, jangan lupa mengarahkan *domain/subdomain* di tab Domains ke aplikasi Anda. (Ubah file `API_URL` di frontend agar mengarah ke domain backend yang telah di-deploy jika diperlukan).
