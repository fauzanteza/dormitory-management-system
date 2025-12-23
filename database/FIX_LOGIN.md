# Fix Login Issue - Panduan Cepat

## Masalah

Login gagal dengan "Invalid credentials" meskipun data sudah diimport.

## Kemungkinan Penyebab

1. Data seed.sql tidak terimport dengan benar
2. Password hash tidak compatible
3. Email salah

## ✅ Solusi Cepat

### Option 1: Buat Test User (TERCEPAT)

**Via phpMyAdmin:**

1. Buka phpMyAdmin → database `dormitory_management`
2. Tab "SQL"
3. Copy-paste isi file `create_test_users.sql`
4. Click "Go"

**Login Credentials:**

- Email: `admin@asrama.com`
- Password: `password123`

Atau:

- Email: `test@student.com`  
- Password: `password123`

### Option 2: Cek Data Existing

**Via phpMyAdmin:**

1. Database `dormitory_management` → tabel `users`
2. Cek apakah ada data
3. Jika kosong → import `seed.sql` lagi

### Option 3: Registrasi User Baru

Gunakan endpoint registrasi:

**POST** `http://localhost:8080/api/auth/register`

Body:

```json
    {
    "name": "Test Admin",
    "email": "newadmin@test.com",
    "password": "test123",
    "role": "admin",
    "phone": "08123456789"
    }
```

Lalu login dengan:

- Email: `newadmin@test.com`
- Password: `test123`

## 🔍 Debugging

### Cek di phpMyAdmin

```sql
SELECT id, name, email, role FROM users;
```

**Hasil yang benar:**

- Harus ada minimal 1 user
- Email harus sesuai
- Role: admin atau student

### Test Password Hash

Password `password123` di-hash jadi:

```
$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW
```

## 🎯 Ringkasan

**Cara Tercepat:**

1. phpMyAdmin → SQL  
2. Jalankan `create_test_users.sql`
3. Login: `admin@asrama.com` / `password123`
