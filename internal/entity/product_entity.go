package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string         `gorm:"size:100;not null"`
	SKU         string         `gorm:"size:50;unique;not null"`
	Price       int            `gorm:"not null;default:0"`
	Stock       int            `gorm:"not null;default:0"`
	Description string         `gorm:"type:text"`
	Star        float64        `gorm:"size:20;default:5.0"`
	Images      []ProductImage `gorm:"foreignKey:ProductID"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null"`    // foreign key
	Category    Category       `gorm:"foreignKey:CategoryID"` // relasi
	SpecialType string         `gorm:"size:50;default:null"`  // contoh: "special1", "special2"
	IsSpecial   bool           `gorm:"default:false"`         // flag khusus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
