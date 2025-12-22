-- Buat database
CREATE DATABASE IF NOT EXISTS dormitory_management;
USE dormitory_management;

-- Tabel users (admin dan mahasiswa)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'student') DEFAULT 'student',
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel rooms
CREATE TABLE rooms (
    id INT PRIMARY KEY AUTO_INCREMENT,
    room_number VARCHAR(10) UNIQUE NOT NULL,
    building VARCHAR(50),
    floor INT,
    capacity INT DEFAULT 2,
    current_occupancy INT DEFAULT 0,
    status ENUM('available', 'occupied', 'maintenance') DEFAULT 'available',
    monthly_rate DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel residents
CREATE TABLE residents (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT UNIQUE NOT NULL,
    room_id INT,
    student_id VARCHAR(20) UNIQUE,
    faculty VARCHAR(100),
    major VARCHAR(100),
    year_of_entry YEAR,
    emergency_contact VARCHAR(20),
    status ENUM('active', 'inactive', 'graduated') DEFAULT 'active',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL
);

-- Tabel payments
CREATE TABLE payments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    resident_id INT NOT NULL,
    room_id INT NOT NULL,
    month DATE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_date DATE,
    status ENUM('paid', 'pending', 'overdue') DEFAULT 'pending',
    payment_method VARCHAR(50),
    receipt_number VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resident_id) REFERENCES residents(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    UNIQUE KEY unique_payment (resident_id, month)
);

-- Tabel repair_requests
CREATE TABLE repair_requests (
    id INT PRIMARY KEY AUTO_INCREMENT,
    resident_id INT NOT NULL,
    room_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    priority ENUM('low', 'medium', 'high') DEFAULT 'medium',
    status ENUM('pending', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
    reported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    assigned_to INT,
    technician_notes TEXT,
    FOREIGN KEY (resident_id) REFERENCES residents(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);

-- Tabel notifications
CREATE TABLE notifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk optimasi query
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_month ON payments(month);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_repair_status ON repair_requests(status);
CREATE INDEX idx_repair_priority ON repair_requests(priority);
