# Sistem Manajemen Asrama

Sistem informasi manajemen asrama mahasiswa dengan antarmuka terpisah untuk Admin dan Mahasiswa.

## 🚀 Fitur Utama

### Admin Features

- **Dashboard** - Statistik sistem lengkap
- **Manajemen Kamar** - CRUD kamar dengan filter
- **Pemesanan** - Approve/reject pemesanan
- **Data Penghuni** - Manajemen penghuni asrama
- **Pembayaran** - Kelola pembayaran & monitoring
- **Perbaikan** - Handle laporan perbaikan
- **Manajemen User** - Kelola akun pengguna

### Student Features

- **Dashboard Personal** - Statistik pribadi
- **Kamar Saya** - Info kamar & teman sekamar
- **Pemesanan** - Pesan kamar tersedia
- **Pembayaran** - Lihat & bayar tagihan
- **Laporan Perbaikan** - Laporkan kerusakan
- **Profil** - Edit profil & ganti password

## 🛠️ Tech Stack

**Backend:**

- Go 1.21+
- Gin Web Framework
- GORM
- JWT Authentication
- MySQL Database

**Frontend:**

- HTML5 + Bootstrap 5
- Vanilla JavaScript
- Font Awesome Icons
- Responsive Design

## 📁 Struktur Proyek

```
dormitory-management-system/
├── backend/
│   ├── config/          # Configuration
│   ├── handlers/        # API Handlers
│   ├── middleware/      # Auth & Rate Limiting
│   ├── models/          # Database Models
│   ├── routes/          # API Routes
│   ├── utils/           # Utilities
│   └── main.go
├── frontend/
│   ├── admin/           # Admin Interface
│   ├── student/         # Student Interface
│   ├── common/          # Shared Resources
│   └── auth/            # Login & Register
├── database/
│   └── schema.sql       # Database Schema
└── README.md
```

## 🔧 Setup & Installation

### Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Modern web browser

### Installation Steps

1. **Clone Repository**

```bash
git clone <repository-url>
cd dormitory-management-system
```

1. **Database Setup**

```sql
CREATE DATABASE dormitory_db;
```

Import schema:

```bash
mysql -u root -p dormitory_db < database/schema.sql
```

1. **Backend Configuration**

Create `.env` file in `backend/` directory:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=dormitory_db
JWT_SECRET=your-secret-key-here
PORT=8080
```

1. **Install Dependencies**

```bash
cd backend
go mod tidy
```

1. **Run Server**

```bash
go run main.go
```

Server will start at `http://localhost:8080`

## 🎯 Usage

### Admin Access

- URL: <http://localhost:8080/login?role=admin>
- Default credentials:
  - Email: `admin@example.com`
  - Password: `admin123`

### Student Access

- URL: <http://localhost:8080/login?role=student>
- Register first or use default:
  - Email: `budi@student.com`
  - Password: `admin123`

## 📡 API Endpoints

### Authentication

- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register (Student)

### Admin API (`/api/admin/*`)

- Dashboard, Rooms, Bookings, Residents
- Payments, Repairs, Users

### Student API (`/api/student/*`)

- Dashboard, My Room, Bookings
- Payments, Repairs, Profile

## 🔐 Security Features

- JWT Token Authentication
- Role-based Access Control
- Rate Limiting (10 req/sec)
- Password Hashing
- Input Validation
- CORS Protection

## 📱 Screenshots

*Coming soon...*

## 🤝 Contributing

1. Fork the project
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 👥 Authors

- Fauzan Teza

## 📞 Support

For support, email: <fauzan@example.com>

---

**Note:** Remember to change default credentials and JWT secret in production!
