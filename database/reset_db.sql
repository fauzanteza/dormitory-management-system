-- Drop dan Recreate Database untuk GORM
DROP DATABASE IF EXISTS dormitory_management;
CREATE DATABASE dormitory_management;
USE dormitory_management;
-- Tabel akan otomatis dibuat oleh GORM AutoMigrate
-- Setelah jalankan server, import seed.sql