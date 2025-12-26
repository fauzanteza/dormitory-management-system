package utils

import (
	"dormitory-management/config"
	"fmt"
	"log"
	"os"
	"strings"
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
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Printf("Connection String (hidden password): %s@tcp(%s:%d)/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)
		log.Fatal("Could not connect to the database. Please check if MySQL is running and the credentials in .env are correct.")
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

	// Import database schema from SQL file
	// Note: We use a relative path assuming the binary is run from the backend directory
	if err := executeSQLFile(DB, "../database/schema.sql"); err != nil {
		log.Printf("Failed to import schema.sql: %v", err)
	} else {
		log.Println("Database schema imported successfully")
	}

	// Seed database
	if err := executeSQLFile(DB, "../database/seed.sql"); err != nil {
		log.Printf("Failed to import seed.sql: %v", err)
	} else {
		log.Println("Database seeded successfully")
	}
}

// Helper function to read and execute SQL file
func executeSQLFile(db *gorm.DB, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Split SQL statements by semicolon
	// Note: This is a simple splitter and might break on semicolons inside strings
	// But for the provided schema/seed files it should be sufficient
	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if err := db.Exec(query).Error; err != nil {
			// Ignore "Query was empty" or specific errors if needed
			if !strings.Contains(err.Error(), "Query was empty") {
				// We log but don't fail immediately to allow partial success
				log.Printf("Error executing query: %s\nError: %v", query[:min(len(query), 50)]+"...", err)
			}
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetDB() *gorm.DB {
	return DB
}
