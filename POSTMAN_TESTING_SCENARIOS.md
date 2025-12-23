# Testing Scenarios - Postman Complete Guide

## Dormitory Management System API

Server: `http://localhost:8080`

---

# 🔐 PART 1: ADMIN WORKFLOW

## Scenario A: Admin Setup & Room Management

### Step 1: Register Admin

**POST** `/api/auth/register`

**Body:**

```json
{
    "name": "Admin Utama",
    "email": "admin@asrama.com",
    "password": "admin123",
    "role": "admin",
    "phone": "081234567890"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "name": "Admin Utama",
    "email": "admin@asrama.com",
    "role": "admin",
    "phone": "081234567890"
}
```

---

### Step 2: Admin Login

**POST** `/api/auth/login`

**Body:**

```json
{
    "email": "admin@asrama.com",
    "password": "admin123"
}
```

**Expected Response (200):**

```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": 1,
        "name": "Admin Utama",
        "email": "admin@asrama.com",
        "role": "admin"
    }
}
```

**⚠️ COPY TOKEN untuk step selanjutnya!**

---

### Step 3: View Dashboard Stats

**GET** `/api/admin/dashboard/stats`

**Headers:**

```
Authorization: Bearer <token_dari_step_2>
```

**Expected Response (200):**

```json
{
    "total_rooms": 0,
    "occupied_rooms": 0,
    "available_rooms": 0,
    "total_residents": 0,
    "pending_payments": 0,
    "pending_bookings": 0,
    "total_revenue": 0,
    "pending_repairs": 0
}
```

---

### Step 4: Create Room - Gedung A

**POST** `/api/admin/rooms`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

```json
{
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1500000,
    "description": "Kamar standar dengan AC dan kamar mandi dalam"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "current_occupancy": 0,
    "status": "available",
    "monthly_rate": 1500000,
    "description": "Kamar standar dengan AC dan kamar mandi dalam"
}
```

---

### Step 5: Create More Rooms

Ulangi Step 4 dengan data berbeda:

**Room A102:**

```json
{
    "room_number": "A102",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1500000,
    "description": "Kamar standar dengan AC"
}
```

**Room B201:**

```json
{
    "room_number": "B201",
    "building": "Gedung B",
    "floor": 2,
    "capacity": 2,
    "monthly_rate": 1600000,
    "description": "Kamar premium dengan balkon"
}
```

**Room C301:**

```json
{
    "room_number": "C301",
    "building": "Gedung C",
    "floor": 3,
    "capacity": 2,
    "monthly_rate": 1700000,
    "description": "Kamar VIP dengan pemandangan"
}
```

---

### Step 6: View All Rooms

**GET** `/api/admin/rooms`

**Headers:**

```
Authorization: Bearer <token>
```

**Expected Response (200):**

```json
[
    {
        "id": 1,
        "room_number": "A101",
        "building": "Gedung A",
        "status": "available",
        "monthly_rate": 1500000
    },
    {
        "id": 2,
        "room_number": "A102",
        ...
    }
]
```

---

### Step 7: View Available Rooms Only

**GET** `/api/admin/rooms/available`

**Headers:**

```
Authorization: Bearer <token>
```

---

### Step 8: Get Specific Room Details

**GET** `/api/admin/rooms/1`

**Headers:**

```
Authorization: Bearer <token>
```

**Expected Response (200):**

```json
{
    "id": 1,
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "current_occupancy": 0,
    "status": "available",
    "monthly_rate": 1500000,
    "description": "Kamar standar dengan AC dan kamar mandi dalam",
    "residents": []
}
```

---

### Step 9: Update Room

**PUT** `/api/admin/rooms/1`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

```json
{
    "room_number": "A101",
    "building": "Gedung A",
    "floor": 1,
    "capacity": 2,
    "monthly_rate": 1550000,
    "status": "available",
    "description": "Kamar standar dengan AC, water heater, dan WiFi"
}
```

---

### Step 10: Set Room to Maintenance

**PUT** `/api/admin/rooms/2`

**Body:**

```json
{
    "status": "maintenance"
}
```

---

## Scenario B: User Management

### Step 11: View All Users

**GET** `/api/admin/users`

**Headers:**

```
Authorization: Bearer <token>
```

---

### Step 12: Create Student User

**POST** `/api/auth/register`

**Body:**

```json
{
    "name": "Budi Santoso",
    "email": "budi@student.com",
    "password": "budi123",
    "role": "student",
    "phone": "081234567891"
}
```

Ulangi untuk student lain:

- Siti Rahayu (<siti@student.com>)
- Ahmad Fadli (<ahmad@student.com>)

---

### Step 13: View Specific User

**GET** `/api/admin/users/2`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 14: Update User Info

**PUT** `/api/admin/users/2`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "name": "Budi Santoso",
    "email": "budi@student.com",
    "phone": "081234567899"
}
```

---

## Scenario C: Resident Management

### Step 15: Create Resident Record

**POST** `/api/admin/residents`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "user_id": 2,
    "room_id": 1,
    "student_id": "STD2024001",
    "faculty": "Fakultas Teknik",
    "major": "Teknik Informatika",
    "year_of_entry": 2024,
    "emergency_contact": "081234560001",
    "status": "active"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "user_id": 2,
    "room_id": 1,
    "student_id": "STD2024001",
    "faculty": "Fakultas Teknik",
    "major": "Teknik Informatika",
    "year_of_entry": 2024,
    "status": "active"
}
```

---

### Step 16: Create More Residents

Ulangi Step 15 untuk:

**Resident 2 (Siti):**

```json
{
    "user_id": 3,
    "room_id": 3,
    "student_id": "STD2024002",
    "faculty": "Fakultas Ekonomi",
    "major": "Manajemen",
    "year_of_entry": 2024,
    "emergency_contact": "081234560002"
}
```

---

### Step 17: View All Residents

**GET** `/api/admin/residents`

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Expected Response:**

```json
[
    {
        "id": 1,
        "student_id": "STD2024001",
        "user": {
            "name": "Budi Santoso",
            "email": "budi@student.com"
        },
        "room": {
            "room_number": "A101",
            "building": "Gedung A"
        }
    }
]
```

---

### Step 18: Get Resident Details

**GET** `/api/admin/residents/1`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 19: Update Resident

**PUT** `/api/admin/residents/1`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "room_id": 2,
    "emergency_contact": "081234560099"
}
```

---

## Scenario D: Payment Management

### Step 20: Create Payment Record

**POST** `/api/admin/payments`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "resident_id": 1,
    "month": "2024-01-01",
    "amount": 1500000,
    "payment_method": "Transfer Bank",
    "receipt_number": "PAY-2024-001",
    "notes": "Pembayaran bulan Januari 2024"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "resident_id": 1,
    "room_id": 1,
    "month": "2024-01-01",
    "amount": 1500000,
    "status": "pending",
    "payment_method": "Transfer Bank",
    "receipt_number": "PAY-2024-001"
}
```

---

### Step 21: View All Payments

**GET** `/api/admin/payments`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 22: Update Payment Status to Paid

**PUT** `/api/admin/payments/1`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "status": "paid",
    "payment_date": "2024-01-05"
}
```

---

### Step 23: Create Pending Payment

**POST** `/api/admin/payments`

**Body:**

```json
{
    "resident_id": 1,
    "month": "2024-02-01",
    "amount": 1500000,
    "notes": "Pembayaran bulan Februari 2024"
}
```

---

### Step 24: Mark Payment as Overdue

**PUT** `/api/admin/payments/2`

**Body:**

```json
{
    "status": "overdue"
}
```

---

## Scenario E: Booking Management

### Step 25: View All Bookings

**GET** `/api/admin/bookings`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 26: Approve Booking

**PUT** `/api/admin/bookings/1`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "status": "approved"
}
```

---

### Step 27: Reject Booking

**PUT** `/api/admin/bookings/2`

**Body:**

```json
{
    "status": "rejected"
}
```

---

## Scenario F: Repair Request Management

### Step 28: View All Repair Requests

**GET** `/api/admin/repairs`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 29: Update Repair Status

**PUT** `/api/admin/repairs/1`

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Body:**

```json
{
    "status": "in_progress",
    "assigned_to": 1,
    "technician_notes": "Teknisi sedang menuju lokasi"
}
```

---

### Step 30: Complete Repair

**PUT** `/api/admin/repairs/1`

**Body:**

```json
{
    "status": "completed",
    "technician_notes": "AC sudah diperbaiki dan berfungsi normal"
}
```

---

### Step 31: Delete Repair Request

**DELETE** `/api/admin/repairs/1`

**Headers:**

```
Authorization: Bearer <admin_token>
```

---

### Step 32: Final Dashboard Check

**GET** `/api/admin/dashboard/stats`

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Expected Response:** Semua angka sudah terisi sesuai data yang dibuat

---

# 👨‍🎓 PART 2: STUDENT WORKFLOW

## Scenario G: Student Registration & Booking

### Step 33: Register as Student

**POST** `/api/auth/register`

**Body:**

```json
{
    "name": "Dewi Lestari",
    "email": "dewi@student.com",
    "password": "dewi123",
    "role": "student",
    "phone": "081234567894"
}
```

---

### Step 34: Student Login

**POST** `/api/auth/login`

**Body:**

```json
{
    "email": "dewi@student.com",
    "password": "dewi123"
}
```

**⚠️ COPY STUDENT TOKEN!**

---

### Step 35: View Available Rooms

**GET** `/api/admin/rooms/available`

**Headers:**

```
Authorization: Bearer <student_token>
```

---

### Step 36: Create Booking Request

**POST** `/api/student/bookings`

**Headers:**

```
Authorization: Bearer <student_token>
Content-Type: application/json
```

**Body:**

```json
{
    "room_id": 4,
    "start_date": "2024-02-01",
    "duration_months": 12,
    "notes": "Ingin kamar di lantai atas dengan pemandangan"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "user_id": 4,
    "room_id": 4,
    "start_date": "2024-02-01",
    "duration_months": 12,
    "status": "pending",
    "notes": "Ingin kamar di lantai atas dengan pemandangan"
}
```

---

### Step 37: View My Bookings

**GET** `/api/student/bookings`

**Headers:**

```
Authorization: Bearer <student_token>
```

---

### Step 38: Cancel Booking

**PUT** `/api/student/bookings/1`

**Headers:**

```
Authorization: Bearer <student_token>
Content-Type: application/json
```

**Body:**

```json
{
    "status": "cancelled"
}
```

---

## Scenario H: Student as Resident

### Step 39: View My Profile

**GET** `/api/auth/profile`

**Headers:**

```
Authorization: Bearer <student_token>
```

**Expected Response:**

```json
{
    "id": 4,
    "name": "Dewi Lestari",
    "email": "dewi@student.com",
    "role": "student",
    "resident": {
        "student_id": "STD2024003",
        "room": {
            "room_number": "C301",
            "building": "Gedung C"
        }
    }
}
```

---

### Step 40: Update My Profile

**PUT** `/api/auth/profile`

**Headers:**

```
Authorization: Bearer <student_token>
Content-Type: application/json
```

**Body:**

```json
{
    "name": "Dewi Lestari S.Kom",
    "phone": "081234567895"
}
```

---

### Step 41: Change Password

**POST** `/api/auth/change-password`

**Headers:**

```
Authorization: Bearer <student_token>
Content-Type: application/json
```

**Body:**

```json
{
    "old_password": "dewi123",
    "new_password": "dewi12345"
}
```

---

## Scenario I: Student Repair Request

### Step 42: Submit Repair Request

**POST** `/api/student/repairs`

**Headers:**

```
Authorization: Bearer <student_token>
Content-Type: application/json
```

**Body:**

```json
{
    "room_id": 4,
    "title": "AC Tidak Dingin",
    "description": "AC di kamar tidak berfungsi dengan baik, udara tidak dingin sejak 2 hari yang lalu",
    "priority": "high"
}
```

**Expected Response (201):**

```json
{
    "id": 1,
    "resident_id": 3,
    "room_id": 4,
    "title": "AC Tidak Dingin",
    "description": "AC di kamar tidak berfungsi dengan baik...",
    "priority": "high",
    "status": "pending"
}
```

---

### Step 43: View My Repair Requests

**GET** `/api/student/repairs`

**Headers:**

```
Authorization: Bearer <student_token>
```

---

### Step 44: Submit Another Repair (Low Priority)

**POST** `/api/student/repairs`

**Body:**

```json
{
    "room_id": 4,
    "title": "Lampu Kamar Mandi Mati",
    "description": "Lampu di kamar mandi tidak menyala",
    "priority": "low"
}
```

---

## Scenario J: Student Payment

### Step 45: View My Payments

**GET** `/api/student/payments`

**Headers:**

```
Authorization: Bearer <student_token>
```

**Expected Response:**

```json
[
    {
        "id": 1,
        "month": "2024-01-01",
        "amount": 1700000,
        "status": "pending",
        "room": {
            "room_number": "C301",
            "monthly_rate": 1700000
        }
    }
]
```

---

### Step 46: View Payment Details

**GET** `/api/student/payments/1`

**Headers:**

```
Authorization: Bearer <student_token>
```

---

# 🧪 PART 3: TESTING SCENARIOS

## Test 1: Full Admin Flow

1. Register Admin (Step 1)
2. Login Admin (Step 2)
3. Create 4 Rooms (Steps 4-5)
4. Create 3 Students (Step 12)
5. Create 2 Residents (Steps 15-16)
6. Create Payments (Step 20)
7. View Dashboard (Step 32)

## Test 2: Full Student Flow

1. Register Student (Step 33)
2. Login Student (Step 34)
3. View Available Rooms (Step 35)
4. Create Booking (Step 36)
5. Submit Repair Request (Step 42)
6. View My Payments (Step 45)

## Test 3: Booking Approval Flow

1. Student creates booking (Step 36)
2. Admin views bookings (Step 25)
3. Admin approves booking (Step 26)
4. Student becomes resident
5. Payment auto-created

## Test 4: Repair Request Flow

1. Student submits repair (Step 42)
2. Admin views repairs (Step 28)
3. Admin assigns technician (Step 29)
4. Admin completes repair (Step 30)

---

# ⚠️ Expected Errors for Testing

## Test Unauthorized Access

**GET** `/api/admin/users`
**Without Authorization header**
→ Expected: 401 Unauthorized

## Test Student Access Admin Endpoint

**GET** `/api/admin/users`
**With student token**
→ Expected: 403 Forbidden

## Test Invalid Login

**POST** `/api/auth/login`
**Wrong password**
→ Expected: 401 Invalid credentials

## Test Duplicate Email

**POST** `/api/auth/register`
**Email yang sudah ada**
→ Expected: 400 Bad Request

---

# ✅ Success Criteria

- [ ] Admin dapat CRUD semua resources
- [ ] Student dapat view dan create booking
- [ ] Student dapat submit repair request
- [ ] Student dapat view payments
- [ ] Authentication berfungsi (token valid)
- [ ] Authorization berfungsi (role-based access)
- [ ] Dashboard stats akurat
- [ ] All endpoints return proper HTTP status codes
