package utils

import (
	"dormitory-management/config"
	"dormitory-management/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(config *config.Config) {
	// Build DSN with proper format for empty password
	var dsn string
	if config.DBPassword == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")

	// Auto Migrate all models to sync with database tables
	// Order: Parent tables first (no dependencies), then child tables
	log.Println("Running auto migration...")

	// Disable foreign key checks temporarily to handle existing tables
	DB.Exec("SET FOREIGN_KEY_CHECKS=0")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Resident{},
		&models.Payment{},
		&models.RepairRequest{},
		&models.Booking{},
	)

	// Re-enable foreign key checks
	DB.Exec("SET FOREIGN_KEY_CHECKS=1")

	if err != nil {
		log.Fatal("Failed to auto migrate database:", err)
	}

	log.Println("Database migration completed successfully")
}

func GetDB() *gorm.DB {
	return DB
}
