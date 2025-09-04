package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	ImageURL  string    `gorm:"size:255;not null"`
	ImageType string    `gorm:"size:50;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
