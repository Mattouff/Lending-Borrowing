package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole defines the role of a user in the platform
type UserRole string

const (
	// RoleUser is a regular user
	RoleUser UserRole = "user"
	// RoleAdmin is an admin user with special privileges
	RoleAdmin UserRole = "admin"
)

// User represents a user in the lending/borrowing platform
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Address   string         `json:"address" gorm:"type:varchar(42);unique;not null"`
	Username  string         `json:"username" gorm:"type:varchar(50);unique"`
	Role      UserRole       `json:"role" gorm:"type:varchar(20);default:'user'"`
	Verified  bool           `json:"verified" gorm:"default:false"`
	Nonce     string         `json:"nonce" gorm:"type:varchar(100)"` // For signature verification
	LastLogin *time.Time     `json:"lastLogin"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
