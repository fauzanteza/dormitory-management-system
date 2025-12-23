# Panduan Postman - Dormitory Management API

## 🚀 Persiapan

### 1. Install Postman

- Download: <https://www.postman.com/downloads/>
- Atau gunakan versi web: <https://web.postman.com/>

### 2. Pastikan Server Running

```bash
cd backend
go run main.go
```

Server harus running di: **<http://localhost:8080>**

---

## 📋 Daftar Endpoint API

### **Authentication (Public - No Token Required)**

#### 1. Register User Baru

**POST** `http://localhost:8080/api/auth/register`

**Headers:**

```
Content-Type: application/json
```

**Body (JSON):**

```json
{
    "name": "Test Admin",
    "email": "admin@test.com",
    "password": "password123",
    "role": "admin",
    "phone": "08123456789"
}
```

**Response Success (201):**

```json
{
    "id": 1,
    "name": "Test Admin",
    "email": "admin@test.com",
    "role": "admin",
    "phone": "08123456789",
    "created_at": "2024-01-01T10:00:00Z"
}
```

#### 2. Login

**POST** `http://localhost:8080/api/auth/login`

**Headers:**

```
Content-Type: application/json
```

**Body (JSON):**

```json
{
    "email": "admin@test.com",
    "password": "password123"
}
```

**Response Success (200):**

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "name": "Test Admin",
        "email": "admin@test.com",
        "role": "admin",
        "phone": "08123456789"
    }
}
```

**⚠️ PENTING:** Copy token dari response untuk digunakan di request selanjutnya!

---

### **Protected Endpoints (Butuh Token)**

Untuk semua endpoint dibawah, tambahkan header:

**Headers:**

```
Content-Type: application/json
Authorization: Bearer <token_dari_login>
```

---

### **Dashboard**

#### 3. Get Dashboard Stats (Admin Only)

**GET** `http://localhost:8080/api/admin/dashboard/stats`

**Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**

```json
{
    "total_rooms": 10,
    "occupied_rooms": 5,
    "available_rooms": 5,
    "total_residents": 8,
    "pending_payments": 3,
    "pending_bookings": 2,
    "total_revenue": 15000000,
    "pending_repairs": 1
}
```

---

### **Rooms Management**

#### 4. Get All Rooms

**GET** `http://localhost:8080/api/admin/rooms`

#### 5. Get Available Rooms

**GET** `http://localhost:8080/api/admin/rooms/available`

#### 6. Get Room by ID

**GET** `http://localhost:8080/api/admin/rooms/1`

#### 7. Create New Room (Admin Only)

**POST** `http://localhost:8080/api/admin/rooms`

**Body:**

```json
{
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1500000,
    "description": "Kamar standar dengan AC"
}
```

#### 8. Update Room

**PUT** `http://localhost:8080/api/admin/rooms/1`

**Body:**

```json
{
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1600000,
    "status": "available",
    "description": "Kamar standar dengan AC dan water heater"
}
```

#### 9. Delete Room

**DELETE** `http://localhost:8080/api/admin/rooms/1`

---

### **Residents Management**

#### 10. Get All Residents

**GET** `http://localhost:8080/api/admin/residents`

#### 11. Get Resident by ID

**GET** `http://localhost:8080/api/admin/residents/1`

#### 12. Create Resident

**POST** `http://localhost:8080/api/admin/residents`

**Body:**

```json
{
    "user_id": 2,
    "room_id": 1,
    "student_id": "STD2024001",
    "faculty": "Teknik",
    "major": "Informatika",
    "year_of_entry": 2024,
    "emergency_contact": "081234567890"
}
```

#### 13. Update Resident

**PUT** `http://localhost:8080/api/admin/residents/1`

#### 14. Delete Resident

**DELETE** `http://localhost:8080/api/admin/residents/1`

---

### **Payments**

#### 15. Get All Payments

**GET** `http://localhost:8080/api/admin/payments`

#### 16. Create Payment

**POST** `http://localhost:8080/api/admin/payments`

**Body:**

```json
{
    "resident_id": 1,
    "month": "2024-01-01",
    "amount": 1500000,
    "payment_method": "Transfer Bank",
    "receipt_number": "PAY-001",
    "notes": "Pembayaran bulan Januari"
}
```

#### 17. Update Payment Status

**PUT** `http://localhost:8080/api/admin/payments/1`

**Body:**

```json
{
    "status": "paid",
    "payment_date": "2024-01-05"
}
```

---

### **Repairs**

#### 18. Get All Repair Requests

**GET** `http://localhost:8080/api/admin/repairs`

#### 19. Create Repair Request

**POST** `http://localhost:8080/api/student/repairs`

**Body:**

```json
{
    "room_id": 1,
    "title": "AC Rusak",
    "description": "AC tidak dingin",
    "priority": "high"
}
```

#### 20. Update Repair Status (Admin)

**PUT** `http://localhost:8080/api/admin/repairs/1`

**Body:**

```json
{
    "status": "completed",
    "technician_notes": "Sudah diperbaiki"
}
```

---

### **Bookings**

#### 21. Get All Bookings

**GET** `http://localhost:8080/api/admin/bookings`

#### 22. Create Booking (Student)

**POST** `http://localhost:8080/api/student/bookings`

**Body:**

```json
{
    "room_id": 1,
    "start_date": "2024-02-01",
    "duration_months": 12,
    "notes": "Ingin kamar di lantai 1"
}
```

#### 23. Approve/Reject Booking (Admin)

**PUT** `http://localhost:8080/api/admin/bookings/1`

**Body:**

```json
{
    "status": "approved"
}
```

---

### **Users Management**

#### 24. Get All Users (Admin)

**GET** `http://localhost:8080/api/admin/users`

#### 25. Get User by ID

**GET** `http://localhost:8080/api/admin/users/1`

#### 26. Update User

**PUT** `http://localhost:8080/api/admin/users/1`

#### 27. Delete User

**DELETE** `http://localhost:8080/api/admin/users/1`

---

## 🎯 Testing Flow Lengkap

### Scenario 1: Register → Login → Create Room

1. **Register Admin:**

   ```
   POST /api/auth/register
   Body: {name, email, password, role: "admin"}
   ```

2. **Login:**

   ```
   POST /api/auth/login
   Body: {email, password}
   → Copy token
   ```

3. **Create Room:**

   ```
   POST /api/admin/rooms
   Headers: Authorization: Bearer <token>
   Body: {room_number, building, floor, capacity, monthly_rate}
   ```

### Scenario 2: Student Booking Flow

1. **Register Student:**

   ```
   POST /api/auth/register
   Body: {name, email, password, role: "student"}
   ```

2. **Login:**

   ```
   POST /api/auth/login
   → Copy token
   ```

3. **View Available Rooms:**

   ```
   GET /api/admin/rooms/available
   Headers: Authorization: Bearer <token>
   ```

4. **Create Booking:**

   ```
   POST /api/student/bookings
   Headers: Authorization: Bearer <token>
   Body: {room_id, start_date, duration_months}
   ```

---

## 📝 Postman Collection

### Import Collection (Optional)

Saya bisa buatkan Postman Collection JSON jika diperlukan. File tersebut bisa langsung di-import ke Postman.

---

## ⚠️ Common Errors

### Error 401 Unauthorized

- **Penyebab:** Token tidak valid atau expired
- **Solusi:** Login ulang untuk dapat token baru

### Error 403 Forbidden

- **Penyebab:** Role tidak sesuai (student akses endpoint admin)
- **Solusi:** Gunakan akun dengan role yang benar

### Error 400 Bad Request

- **Penyebab:** Format body JSON salah
- **Solusi:** Cek format JSON di body

### Error 500 Internal Server Error

- **Penyebab:** Error di server
- **Solusi:** Cek log server di terminal

---

## 🔍 Tips Postman

1. **Gunakan Environment Variables:**
   - `base_url`: `http://localhost:8080`
   - `token`: Copy dari response login
   - `admin_token`: Token admin
   - `student_token`: Token student

2. **Gunakan Pre-request Scripts:**

   ```javascript
   pm.environment.set("token", pm.response.json().token);
   ```

3. **Simpan Request di Collection:**
   - Buat folder untuk setiap module
   - Auth, Rooms, Residents, dll

---

## ✅ Quick Test

**Test koneksi API:**

```
GET http://localhost:8080/health
```

**Response:**

```json
{
    "status": "ok",
    "message": "Server is running"
}
```
