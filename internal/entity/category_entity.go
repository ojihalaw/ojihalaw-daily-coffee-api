package entity

import (
	"time"

	"github.com/google/uuid"
)

// Implement Searchable
func (Category) SearchFields() []string {
	return []string{"name"}
}

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"size:100;not null;unique"`
	Slug      string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
