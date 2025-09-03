package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	UserName    string    `gorm:"size:50;not null"`
	Email       string    `gorm:"size:50;unique;not null"`
	Password    string    `gorm:"not null"`
	PhoneNumber string    `gorm:"column:phone_number;size:20"`
	Role        string    `gorm:"size:20;default:customer"`
	Status      string    `gorm:"size:20;default:active"`
	Address     string    `gorm:"size:255"`
	Addresses   string    `gorm:"size:255"`
	City        string    `gorm:"size:100"`
	PostalCode  string    `gorm:"size:10"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (u *Customer) TableName() string {
	return "customers"
}
