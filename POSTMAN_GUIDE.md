# Panduan Lengkap Testing dengan Postman

Panduan ini mencakup langkah-langkah detail untuk menguji **setiap fitur** dalam Dormitory Management System menggunakan koleksi Postman yang telah disediakan.

## Persiapan Awal

1. **Start Server**: Pastikan backend berjalan (`go run main.go` di folder backend).
2. **Import Collection**: Import file `dormitory_api.postman_collection.json` ke Postman.
3. **Environment**: Pastikan variable `base_url` di set ke `http://localhost:8080`.

---

## Skenario 1: Otentikasi (Auth)

Setiap fitur membutuhkan Login. Token akan tersimpan otomatis.

1. **Register Account Baru**
    * Buka folder **Auth** > **Register**.
    * Klik **Body**. Ubah data JSON jika perlu (msl: username, email).
    * Klik **Send**. Pastikan response `201 Created` atau `200 OK`.
2. **Login Admin** (Untuk akses fitur Admin)
    * Buka **Auth** > **Login**.
    * Gunakan credential default:
        * User: `admin`
        * Pass: `admin123`
    * Klik **Send**. Token akan otomatis tersimpan.
3. **Login Student** (Untuk akses fitur Student)
    * Gunakan akun yang baru diregister di langkah 1.
    * Klik **Send**. Token akan otomatis terupdate.

---

## Skenario 2: Manajemen Kamar (Admin Only)

Pastikan Anda Login sebagai **Admin** sebelum menjalankan ini.

1. **Melihat Semua Kamar**
    * **Admin** > **Rooms** > **Get All Rooms**.
    * Klik **Send**. Daftar kamar akan muncul.
2. **Menambah Kamar Baru**
    * **Admin** > **Rooms** > **Create Room**.
    * Klik **Body**. Pastikan semua **field wajib** terisi dengan benar di JSON:
        * `room_number`: Nomor kamar (contoh: "101")
        * `building`: Nama gedung (contoh: "Gedung A")
        * `floor`: Lantai (integer, contoh: 1)
        * `capacity`: Kapasitas orang (integer, contoh: 2)
        * `monthly_rate`: Harga per bulan (float/number, contoh: 1500000)
    * Klik **Send**.
3. **Update Kamar**
    * **Admin** > **Rooms** > **Update Room**.
    * Ubah value `:id` di URL (misal: `1`).
    * Di **Body**, masukkan hanya field yang ingin diubah (misal: `monthly_rate`). Field harus sesuai (misal: `room_number`, `building`, `floor`, dll).
    * Klik **Send**.
4. **Hapus Kamar**
    * **Admin** > **Rooms** > **Delete Room**.
    * Set `:id` kamar yang ingin dihapus.
    * Klik **Send**.

---

## Skenario 3: Alur Booking (Student & Admin)

Simulasi mahasiswa memesan kamar dan admin menyetujuinya.

1. **Student: Melihat Kamar Tersedia**
    * **Login sebagai Student**.
    * **Student** > **Get Available Rooms**.
    * Pilih ID kamar yang tersedia.
2. **Student: Membuat Booking**
    * **Student** > **Create Booking**.
    * Isi `room_id` dengan ID kamar dari langkah sebelumnya.
    * Klik **Send**. Booking akan berstatus `pending`.
3. **Admin: Approve Booking**
    * **Login sebagai Admin**.
    * **Admin** > **Bookings** > **Get All Bookings**.
    * Cari ID booking yang baru dibuat.
    * **Admin** > **Bookings** > **Update Booking Status**.
    * Set `:id` booking tersebut.
    * Isi Body dengan status `approved`.
    * Klik **Send**. Mahasiswa resmi menjadi penghuni (resident).

---

## Skenario 4: Pembayaran (Finance Flow)

Admin membuat tagihan, mahasiswa melihat, dan admin konfirmasi bayar.

1. **Admin: Buat Tagihan (Invoice)**
    * **Admin** > **Payments** > **Create Payment**.
    * Isi `resident_id` (ID penghuni dari hasil booking tadi), jumlah, dan tanggal.
    * Klik **Send**. Status awal: `pending`.
2. **Student: Cek Tagihan**
    * **Login sebagai Student**.
    * **Student** > **Payments** > **Cek My Payments** (Gunakan endpoint My Payments).
3. **Admin: Konfirmasi Pembayaran**
    * **Login sebagai Admin**.
    * **Admin** > **Payments** > **Update Payment Status**.
    * Set `:id` payment.
    * Isi Body dengan status `paid`.
    * Klik **Send**.

---

## Skenario 5: Perbaikan (Maintenance Flow)

Mahasiswa lapor kerusakan, admin update progres.

1. **Student: Lapor Kerusakan**
    * **Login sebagai Student**.
    * **Student** > **Create Repair Request** (di dalam folder Repairs).
    * Isi keluhan (misal: "Kran bocor").
    * Klik **Send**.
2. **Admin: Cek Laporan**
    * **Login sebagai Admin**.
    * **Admin** > **Repairs** > **Get Repair Requests**.
    * Lihat daftar laporan.

---

## Tips Tambahan

* **Error 401 Unauthorized**: Berarti Token habis atau belum login. Lakukan **Login** ulang.
* **Error Validation**: Jika muncul error seperti `Field validation for ... failed`, periksa kembali nama field di JSON Body. Pastikan `room_number`, `building`, `floor` ada pada Create Room.
