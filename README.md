# Sistem Informasi Manajemen Asrama Mahasiswa

Sistem terintegrasi untuk pengelolaan asrama mahasiswa yang dibangun dengan Go, MySQL, dan HTML/CSS/JS.

## Fitur Utama
1. **Manajemen Kamar**
   - Status kamar (tersedia, terisi, maintenance)
   - Kapasitas dan okupansi
   - Harga sewa bulanan
2. **Manajemen Penghuni**
   - Data mahasiswa
   - Penempatan kamar
   - Status keaktifan
3. **Sistem Pembayaran**
   - Pembayaran bulanan
   - Status pembayaran (lunas, pending, terlambat)
   - Laporan bulanan
4. **Sistem Perbaikan**
   - Pengajuan perbaikan
   - Prioritas dan status
   - Penugasan teknisi

## Teknologi Stack
- **Backend**: Go + Gin Framework
- **Database**: MySQL
- **ORM**: GORM
- **Authentication**: JWT
- **Frontend**: HTML5, CSS3, JavaScript
- **Charts**: Chart.js
- **Tables**: DataTables

## Instalasi dan Setup

### Prerequisites
1. Go 1.21 atau versi terbaru
2. MySQL 8.0 atau versi terbaru
3. Git

### Langkah Instalasi

#### 1. Clone Repository
```bash
git clone [repository-url]
cd dormitory-management-system
```

#### 2. Setup Database
```bash
# Masuk ke MySQL sebagai root
mysql -u root -p

# Jalankan script database
source database/schema.sql
source database/seed.sql
```

Atau gunakan script setup:
```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

#### 3. Konfigurasi Environment
```bash
cd backend
cp .env.example .env
# Edit file .env dengan konfigurasi database Anda
```

#### 4. Install Dependencies
```bash
cd backend
go mod tidy
```

#### 5. Jalankan Aplikasi
```bash
cd backend
go run main.go
```
Aplikasi akan berjalan di http://localhost:8080

### Akses Default
**Admin**
Email: admin@asrama.com
Password: admin123

**Mahasiswa**
Email: budi@student.com
Password: admin123

## API Endpoints

### Autentikasi
- `POST /api/auth/login` - Login user
- `POST /api/auth/register` - Registrasi user baru

### Kamar
- `GET /api/rooms` - Get semua kamar
- `GET /api/rooms/available` - Get kamar tersedia
- `POST /api/rooms` - Tambah kamar baru (admin only)
- `PUT /api/rooms/:id` - Update kamar (admin only)
- `DELETE /api/rooms/:id` - Hapus kamar (admin only)

### Pembayaran
- `GET /api/payments` - Get semua pembayaran
- `POST /api/payments` - Buat pembayaran baru (admin only)
- `PUT /api/payments/:id/status` - Update status pembayaran (admin only)
- `GET /api/payments/report` - Get laporan bulanan

### Perbaikan
- `GET /api/repairs` - Get semua permintaan perbaikan
- `POST /api/repairs` - Buat permintaan perbaikan baru
- `PUT /api/repairs/:id/status` - Update status perbaikan (admin only)

### Dashboard
- `GET /api/dashboard/stats` - Get statistik dashboard
