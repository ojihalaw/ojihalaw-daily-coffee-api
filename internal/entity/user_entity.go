package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
	RoleFinance Role = "finance"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	UserName    string    `gorm:"size:50;not null"`
	Email       string    `gorm:"size:50;unique;not null"`
	Password    string    `gorm:"not null"`
	PhoneNumber string    `gorm:"column:phone_number;size:20"`
	Role        Role      `gorm:"type:varchar(20);not null;default:'user'"`
	Status      string    `gorm:"size:20;default:active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
