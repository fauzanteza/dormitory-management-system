# Cara Import Seed Data

## ✅ Status Database

- Database: `dormitory_management` sudah ada
- Tables: 6 tabel sudah dibuat oleh GORM
- Server: Running di port 8080

## 📝 Import seed.sql

**TIDAK PERLU** membuat database baru atau jalankan schema.sql!

Tabel sudah compatible, tinggal import data:

### Option 1: Via phpMyAdmin (MUDAH)

1. Buka **phpMyAdmin** (biasanya <http://localhost/phpmyadmin>)
2. Klik database **dormitory_management** di sidebar kiri
3. Klik tab **"Import"** di menu atas
4. Klik **"Choose File"** → pilih `seed.sql` di folder database
5. Scroll ke bawah, klik **"Go"**

### Option 2: Via MySQL Command Line

Buka terminal di folder database:

```bash
mysql -u root dormitory_management < seed.sql
```

Atau di PowerShell:

```powershell
Get-Content seed.sql | mysql -u root dormitory_management
```

## 📊 Data yang Akan Diimport

- 4 Users (1 admin, 3 students)
- 6 Rooms (Gedung A, B, C)
- 2 Residents
- 4 Payments
- 2 Repair Requests

## 🔑 Login Setelah Import

- Admin: `admin@example.com` / `admin123`
- Student: `student@example.com` / password dari seed

## ✅ Verifikasi

Setelah import, cek di phpMyAdmin:

- Tabel users → harus ada 4 baris
- Tabel rooms → harus ada 6 baris
