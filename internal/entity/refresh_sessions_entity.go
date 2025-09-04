package entity

import (
	"time"

	"github.com/google/uuid"
)

// Kenapa simpan refresh_session?
// Untuk token rotation (revoke lama saat refresh), force logout, audit, dan memitigasi theft.
type RefreshSession struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	JTI       string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	Revoked   bool      `gorm:"default:false"`
	IP        string    `gorm:"size:64"`
	UserAgent string    `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
