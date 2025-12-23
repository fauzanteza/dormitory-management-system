# Database Setup Guide

## Status Saat Ini ✅

Database **sudah berjalan** dengan GORM AutoMigrate!

- **Database:** dormitory_management
- **Total Tabel:** 6 (users, rooms, residents, payments, repair_requests, bookings)
- **Server:** Running di port 8080

## Cara Import Seed Data

### Option 1: Via MySQL Command (Jika ada MySQL CLI)

```bash
mysql -u root -p dormitory_management < database/seed_gorm.sql
```

### Option 2: Via phpMyAdmin / MySQL Workbench

1. Buka phpMyAdmin atau MySQL Workbench
2. Pilih database `dormitory_management`
3. Import file `database/seed_gorm.sql`

### Option 3: Via Aplikasi (Otomatis)

Server akan import seed data secara otomatis saat pertama kali jalan (jika tabel kosong)

## File Database

### File Baru (Compatible dengan GORM)

- `schema_gorm.sql` - Dokumentasi struktur tabel yang dibuat GORM
- `seed_gorm.sql` - Data sample untuk testing (RECOMMENDED)

### File Lama (Untuk Referensi)

- `schema.sql` - Schema asli dengan INT (tidak compatible)
- `seed.sql` - Seed data asli (perlu disesuaikan)

## User Accounts Setelah Import

### Admin

- Email: `admin@asrama.com`
- Password: `password123`

### Students

- Email: `fauzan@student.com` | Password: `password123`
- Email: `budi@student.com` | Password: `password123`
- Email: `siti@student.com` | Password: `password123`
- Email: `dewi@student.com` | Password: `password123`

## Sample Data yang Akan Diimport

- **5 Users** (1 admin, 4 students)
- **8 Rooms** (Gedung A, B, C)
- **2 Residents** (penghuni aktif)
- **2 Bookings** (1 pending, 1 approved)
- **4 Payments** (2 paid, 1 pending, 1 overdue)
- **3 Repair Requests** (berbagai status)

## Verifikasi Database

Jalankan script verifikasi:

```bash
cd backend
go run verify_db.go
```

## Troubleshooting

### Jika ada error saat import seed

1. Pastikan database `dormitory_management` sudah ada
2. Pastikan semua tabel sudah dibuat (6 tabel)
3. Cek apakah server sedang running (stop dulu sebelum import)

### Reset Database

Jika ingin mulai dari awal:

```sql
DROP DATABASE dormitory_management;
CREATE DATABASE dormitory_management;
```

Lalu jalankan server lagi untuk auto-create tables, kemudian import seed.
