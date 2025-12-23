-- Fix Login: Create Test Users dengan Password yang Benar
USE dormitory_management;
-- Hapus user lama kalau ada
DELETE FROM users
WHERE email IN ('admin@asrama.com', 'test@student.com');
-- Password: "password123" (di-hash dengan bcrypt cost 10)
-- Hash ini 100% kompatibel dengan bcrypt Go
INSERT INTO users (
        name,
        email,
        password,
        role,
        phone,
        created_at,
        updated_at
    )
VALUES (
        'Admin Test',
        'admin@asrama.com',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW',
        'admin',
        '081234567890',
        NOW(),
        NOW()
    ),
    (
        'Student Test',
        'test@student.com',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBoSLP4HGdKjYW',
        'student',
        '081234567891',
        NOW(),
        NOW()
    );
-- Verifikasi
SELECT id,
    name,
    email,
    role,
    SUBSTRING(password, 1, 20) as password_preview
FROM users
WHERE email IN ('admin@asrama.com', 'test@student.com');
-- HASIL: 
-- Harus ada 2 users dengan password yang dimulai dengan $2a$10$