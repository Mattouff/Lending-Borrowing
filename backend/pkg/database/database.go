package database

import (
	"fmt"
	"os"
	"time"

	"github.com/Mattouff/Lending-Borrowing/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		host := os.Getenv("PGHOST")
		user := os.Getenv("PGUSER")
		password := os.Getenv("PGPASSWORD")
		dbname := os.Getenv("PGDATABASE")

		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		databaseURL = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			host, user, password, dbname, port,
		)
	}

	maxRetries := 5
	retryInterval := time.Second * 3

	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Connected to database")
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	// Auto-migrate all models
	return db.AutoMigrate(
		&models.User{},
	)
}
