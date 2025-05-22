package database

import (
	"log"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"gorm.io/gorm"
)

// MigrateDB runs database migrations to create or update tables
func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// List all models that should be migrated
	// Order matters for foreign key dependencies
	err := db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Position{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

// SeedDB seeds the database with initial data if needed
func SeedDB(db *gorm.DB) error {
	// Check if admin user exists
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)

	if count == 0 {
		// Create admin user if none exists
		admin := models.User{
			Address:  "0x1234567890123456789012345678901234567890", // Should be replaced with actual admin address
			Username: "admin",
			Role:     models.RoleAdmin,
			Verified: true,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to seed admin user: %v", err)
			return err
		}
		log.Println("Admin user seeded successfully")
	}

	return nil
}

// ResetDB drops all tables and recreates them
// WARNING: This will delete all data. REMOVE BEFORE PRODUCTION!
func ResetDB(db *gorm.DB) error {
	log.Println("WARNING: Resetting database (all data will be lost)...")

	err := db.Migrator().DropTable(
		&models.Position{},
		&models.Transaction{},
		&models.User{},
	)

	if err != nil {
		log.Printf("Failed to drop tables: %v", err)
		return err
	}

	return MigrateDB(db)
}
