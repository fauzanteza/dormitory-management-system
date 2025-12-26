USE dormitory_management;
-- Insert admin user (password: admin123)
INSERT IGNORE INTO users (name, email, password, role, phone)
VALUES (
        'fauzan',
        'fauzan@example.com',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW',
        'admin',
        '081234567890'
    ),
    (
        'fauzan2',
        'fauzan2@example.com',
        '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'student',
        '081234567891'
    ),
    (
        'Budi Santoso',
        'budi@student.com',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW',
        'student',
        '081234567891'
    ),
    (
        'Siti Rahayu',
        'siti@student.com',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW',
        'student',
        '081234567892'
    );
-- Insert rooms
INSERT INTO rooms (
        room_number,
        building,
        floor,
        capacity,
        monthly_rate,
        status
    )
VALUES (
        'A101',
        'Gedung A',
        1,
        2,
        1500000.00,
        'available'
    ),
    ('A102', 'Gedung A', 1, 2, 1500000.00, 'occupied'),
    (
        'A103',
        'Gedung A',
        1,
        2,
        1500000.00,
        'available'
    ),
    ('B201', 'Gedung B', 2, 2, 1600000.00, 'occupied'),
    (
        'B202',
        'Gedung B',
        2,
        2,
        1600000.00,
        'maintenance'
    ),
    (
        'C301',
        'Gedung C',
        3,
        2,
        1700000.00,
        'available'
    );
-- Insert residents
INSERT INTO residents (
        user_id,
        room_id,
        student_id,
        faculty,
        major,
        year_of_entry
    )
VALUES (2, 2, '20230001', 'Teknik', 'Informatika', 2023),
    (3, 4, '20230002', 'Ekonomi', 'Manajemen', 2023);
-- Update occupancy
UPDATE rooms
SET current_occupancy = 1
WHERE id IN (2, 4);
-- Insert sample payments
INSERT INTO payments (
        resident_id,
        room_id,
        month,
        amount,
        status,
        payment_date
    )
VALUES (
        1,
        2,
        '2024-01-01',
        1500000.00,
        'paid',
        '2024-01-05'
    ),
    (1, 2, '2024-02-01', 1500000.00, 'pending', NULL),
    (
        2,
        4,
        '2024-01-01',
        1600000.00,
        'paid',
        '2024-01-10'
    ),
    (2, 4, '2024-02-01', 1600000.00, 'overdue', NULL);
-- Insert sample repair requests
INSERT INTO repair_requests (
        resident_id,
        room_id,
        title,
        description,
        priority,
        status
    )
VALUES (
        1,
        2,
        'AC Tidak Dingin',
        'AC di kamar tidak berfungsi dengan baik',
        'high',
        'pending'
    ),
    (
        2,
        4,
        'Keran Bocor',
        'Keran wastafel di kamar mandi bocor',
        'medium',
        'in_progress'
    );