# Panduan Lengkap & Mudah: Uji Coba Full Fitur (Student & Admin)

Panduan ini dirancang agar Anda bisa mencoba **semua fitur** aplikasi dari A sampai Z menggunakan Postman. Kita akan menggunakan skenario nyata:

1. **Admin** menyiapkan kamar.
2. **Student (fauzan2)** mendaftar dan menyewa kamar.
3. **Admin** menyetujui dan menagih pembayaran.
4. **Student** lapor kerusakan.
5. **Admin** memproses laporan.

---

## 🟢 Persiapan (Wajib)

1. Jalankan server backend:

   ```bash
   go run main.go
   ```

2. Pastikan Postman sudah terbuka.

---

## BAGIAN 1: Persiapan Admin (Siapkan Kamar)

Sebelum mahasiswa bisa daftar, harus ada kamar dulu.

### 1. Login Admin

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/auth/login`
- **Body** (JSON):

  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```

- **Response**: Copy **Token** (teks panjang dimulai `eyJ...`).
- **Catatan**: Kita sebut ini **[TOKEN ADMIN]**.

### 2. Buat Kamar Baru

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/admin/rooms`
- **Auth**: Pilih Type **Bearer Token** -> Paste **[TOKEN ADMIN]**.
- **Body** (JSON):

  ```json
  {
    "room_number": "A-101",
    "building": "Tower A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1500000,
    "description": "Kamar AC dengan jendela"
  }
  ```

- **Response**: Perhatikan `id` kamar yang baru dibuat (misal: `1`). Kita sebut ini **[ID KAMAR]**.

---

## BAGIAN 2: Student Journey (fauzan2)

Sekarang kita berperan sebagai mahasiswa baru.

### 3. Register Akun (fauzan2)

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/auth/register`
- **Body** (JSON):

  ```json
  {
    "name": "fauzan2",
    "email": "fauzan2@example.com",
    "password": "fauzan123",
    "role": "student",
    "phone": "0812345678"
  }
  ```

- **Send**: Harusnya berhasil (Status 201).

### 4. Login Student

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/auth/login`
- **Body** (JSON):

  ```json
  {
    "email": "fauzan2@example.com",
    "password": "fauzan123"
  }
  ```

- **Response**: Copy **Token** baru ini. Kita sebut ini **[TOKEN STUDENT]**.
- _PENTING: Sekarang ganti semua Auth di Postman pakai token ini._

### 5. Lihat Profile

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/student/profile`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Send**: Muncul data fauzan2.

### 6. Cari Kamar Kosong

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/student/rooms/available`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Send**: Anda akan lihat kamar "A-101" yang dibuat Admin tadi. Pastikan Anda ingat **ID**-nya (misal: `1`).

### 7. Booking Kamar (Sewa)

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/student/bookings`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Body** (JSON):

  ```json
  {
    "room_id": 1,
    "start_date": "2024-01-01",
    "end_date": "2024-12-31"
  }
  ```

  _(Ganti `room_id` dengan ID kamar yang nyata jika bukan 1)_

- **Send**: Berhasil booking dengan status `pending`.

### 8. Cek Bookingan Saya

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/student/bookings`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Send**: Muncul list booking Anda.

---

## BAGIAN 3: Admin Action (Approve & Tagih)

Kembali jadi Admin untuk menyetujui sewa.

### 9. Approve Booking

- **Method**: `PUT`
- **URL**: `http://localhost:8080/api/admin/bookings/1/status`
  _(Ganti angka `1` di URL dengan ID Booking jika berbeda)_
- **Auth**: Paste **[TOKEN ADMIN]**.
- **Body** (JSON):

  ```json
  {
    "status": "approved"
  }
  ```

- **Send**: Status berubah jadi `approved`. Sistem otomatis membuat data penghuni (Resident).

### 10. Cari Data Penghuni (Untuk Buat Tagihan)

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/admin/residents`
- **Auth**: Paste **[TOKEN ADMIN]**.
- **Send**: Cari data `fauzan2`. Ambil **ID** resident-nya (misal: `1`). Kita sebut **[ID RESIDENT]**.

### 11. Buat Tagihan (Invoice)

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/admin/payments`
- **Auth**: Paste **[TOKEN ADMIN]**.
- **Body** (JSON):

  ```json
  {
    "resident_id": 1,
    "amount": 1500000,
    "month": "2024-01",
    "notes": "Pembayaran Bulan Januari"
  }
  ```

  _(Ganti `resident_id` dengan ID Resident yang asli)_

- **Send**: Tagihan dibuat dengan status `pending`. Copy **ID Payment** (misal: `1`).

---

## BAGIAN 4: Student (Bayar & Lapor Masalah)

Balik lagi jadi `fauzan2`.

### 12. Cek Tagihan Saya

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/student/my-payments`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Send**: Muncul tagihan Rp 1.500.000.

### 13. Lapor Kerusakan (Komplain)

Anggap ada kran bocor.

- **Method**: `POST`
- **URL**: `http://localhost:8080/api/student/repairs`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Body** (JSON):

  ```json
  {
    "room_id": 1,
    "title": "Kran Air Bocor",
    "description": "Air menetes terus di kamar mandi",
    "priority": "high"
  }
  ```

- **Send**: Laporan terkirim.

### 14. Cek Laporan Saya

- **Method**: `GET`
- **URL**: `http://localhost:8080/api/student/my-repairs`
- **Auth**: Paste **[TOKEN STUDENT]**.
- **Send**: Muncul status `pending`.

---

## BAGIAN 5: Admin (Selesaikan Masalah)

Admin memverifikasi pembayaran dan memperbaiki kerusakan.

### 15. Konfirmasi Pembayaran

Admin menandai tagihan sudah dibayar.

- **Method**: `PUT`
- **URL**: `http://localhost:8080/api/admin/payments/1/status`
  _(Ganti `1` dengan ID Payment)_
- **Auth**: Paste **[TOKEN ADMIN]**.
- **Body** (JSON):

  ```json
  {
    "status": "paid"
  }
  ```

- **Send**: Status berubah jadi `paid`.

### 16. Update Status Perbaikan

Teknisi sudah memperbaiki kran.

- **Method**: `PUT`
- **URL**: `http://localhost:8080/api/admin/repairs/1/status`
  _(Ganti `1` dengan ID Repair)_
- **Auth**: Paste **[TOKEN ADMIN]**.
- **Body** (JSON):

  ```json
  {
    "status": "completed"
  }
  ```

- **Send**: Masalah selesai!

---

## Rangkuman

Selamat! Anda sudah mensimulasikan siklus penuh aplikasi:

1. ✅ Admin siapin kamar.
2. ✅ Student daftar & booking.
3. ✅ Admin setujui masuk.
4. ✅ Admin tagih bayaran.
5. ✅ Student lapor rusak.
6. ✅ Admin bereskan semua.
