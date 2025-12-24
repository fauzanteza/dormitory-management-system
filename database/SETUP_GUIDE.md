# Panduan Lengkap: Reset Database & Setup

## ⚠️ Masalah

Error: `Cannot change column 'user_id'` - Tabel dibuat dengan INT, GORM butuh BIGINT UNSIGNED.

## ✅ Solusi: 3 Langkah Mudah

### Step 1: Drop Database via phpMyAdmin

1. **Buka phpMyAdmin**
   - URL: <http://localhost/phpmyadmin>
   - Atau dari XAMPP Control Panel → klik "Admin" pada MySQL

2. **Drop Database**
   - Klik database **`dormitory_management`** di sidebar kiri
   - Klik tab **"Operations"** di menu atas
   - Scroll ke bawah, cari bagian **"Remove database"**
   - Klik tombol **"Drop the database (DROP)"** (warna merah)
   - **Konfirmasi** dengan klik "OK"

3. **Buat Database Baru**
   - Klik tab **"Databases"** di menu atas
   - Tulis nama: `dormitory_management`
   - Collation: `utf8mb4_general_ci`
   - Klik **"Create"**

**Atau via SQL Tab:**

```sql
DROP DATABASE IF EXISTS dormitory_management;
CREATE DATABASE dormitory_management;
```

---

### Step 2: Jalankan Server (Auto-Create Tables)

Di terminal/PowerShell:

```bash
cd backend
go run main.go
```

**Output yang benar:**

2025/12/24 01:XX:XX Database connected successfully
2025/12/24 01:XX:XX Running auto migration...
2025/12/24 01:XX:XX Database migration completed successfully
[GIN-debug] Listening and serving HTTP on :8080

✅ Jika muncul "Database migration completed successfully" → **Berhasil!**

**Biarkan server tetap running** untuk sekarang.

---

### Step 3: Import Seed Data

#### Buka Terminal/PowerShell Baru (jangan tutup yang server)

**Opsi A - Via phpMyAdmin (MUDAH):**

1. Buka phpMyAdmin
2. Klik database `dormitory_management`
3. Tab **"Import"**
4. **"Choose File"** → pilih `database/seed.sql`
5. Klik **"Go"** (scroll ke bawah)

**Opsi B - Via SQL Query:**

1. phpMyAdmin → database `dormitory_management`
2. Tab **"SQL"**
3. Copy-paste isi file `seed.sql`
4. Klik **"Go"**

**Import berhasil jika:**

- Muncul pesan sukses hijau
- Tabel `users` ada 4 baris
- Tabel `rooms` ada 6 baris

---

### Step 4: Restart Server

1. **Stop server** yang running (di terminal pertama: **Ctrl+C**)
2. **Jalankan lagi:**

   ```bash
   go run main.go
   ```

---

## ✅ Verifikasi Berhasil

### Cek Database

- phpMyAdmin → `dormitory_management` → harus ada **6 tabel**
- Tabel `users` → **4 rows** (1 admin, 3 students)
- Tabel `rooms` → **6 rows**

### Cek Server

[GIN-debug] Listening and serving HTTP on :8080

### Test Login

- Buka frontend: `frontend/auth/login.html`
- Login:
  - Email: `admin@example.com`
  - Password: `admin123`

---

## 🎯 Ringkasan Singkat

1. **phpMyAdmin** → Drop database → Create database
2. **Terminal 1:** `go run main.go` → Tabel auto-created
3. **phpMyAdmin** → Import `seed.sql`
4. **Terminal 1:** Ctrl+C → `go run main.go`

**Selesai!** 🚀
