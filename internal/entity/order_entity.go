package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	OrderStatusPending = "pending"
	OrderStatusPaid    = "paid"
	OrderStatusFailed  = "failed"
	OrderStatusExpired = "expired"
)

type Order struct {
	ID            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID    `gorm:"type:uuid;not null"`                 // siapa yang order
	InvoiceNumber string       `gorm:"size:50;unique;not null"`            // kode unik, misal: INV-20250908-0001
	Status        string       `gorm:"size:20;not null;default:'pending'"` // pending, paid, failed, expired
	Amount        int64        `gorm:"not null"`                           // total harga
	PaymentMethod string       `gorm:"size:50"`                            // ex: bank_transfer
	PaymentType   string       `gorm:"size:50"`                            // ex: bca, gopay, shopeepay
	TransactionID string       `gorm:"size:100"`                           // dari Midtrans
	RedirectURL   string       `gorm:"size:255"`                           // kalau pakai Snap
	ExpiredAt     *time.Time   `gorm:"default:null"`
	OrderItems    []OrderItem  `gorm:"foreignKey:OrderID"`
	PaymentLogs   []PaymentLog `gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductName string    `gorm:"size:100;not null"` // disimpan biar invoice tetap valid kalau product berubah
	Qty         int       `gorm:"not null"`
	Price       int64     `gorm:"not null"` // harga per unit saat order
	Subtotal    int64     `gorm:"not null"` // qty * price
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PaymentLog struct {
	ID                  uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID             uuid.UUID      `gorm:"type:uuid;not null"`
	MidtransTransaction string         `gorm:"size:100;not null"`
	Status              string         `gorm:"size:20;not null"`
	RawResponse         datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt           time.Time
}
