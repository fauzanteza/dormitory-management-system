#!/bin/bash

echo "Setting up Dormitory Management System..."

# Install MySQL jika belum ada
if ! command -v mysql &> /dev/null; then
    echo "MySQL tidak ditemukan. Silakan install MySQL terlebih dahulu."
    exit 1
fi

# Buat database dan tabel
echo "Membuat database dan tabel..."
mysql -u root -p < database/schema.sql

echo "Menambahkan data awal..."
mysql -u root -p dormitory_management < database/seed.sql

echo "Setup database selesai!"
