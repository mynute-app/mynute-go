package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
	StatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	BaseModel
	Amount decimal.Decimal `gorm:"type:numeric(12,2);not null"`
	// ISO 4217 currency code (e.g., "USD", "EUR")
	Currency string `gorm:"type:varchar(3);not null;index"`

	// Status
	Status PaymentStatus `gorm:"type:varchar(20);not null;index;default:'PENDING'"`

	// Payment Method Details (adjust as needed)
	PaymentMethod string `gorm:"type:varchar(50);index"`        // e.g., "CREDIT_CARD", "PAYPAL", "BANK_TRANSFER"
	TransactionID string `gorm:"type:varchar(100);uniqueIndex"` // External transaction ID from payment provider (often unique)
	Provider      string `gorm:"type:varchar(50);index"`        // e.g., "Stripe", "PayPal"

	// References to other models (Foreign Keys)
	// Use pointers (*uint) if the relationship is optional (nullable foreign key)
	UserID  *uint `gorm:"index"` // Foreign key to a User model (if applicable)
	OrderID *uint `gorm:"index"` // Foreign key to an Order model (if applicable)
	// If UserID/OrderID are non-nullable, use `uint` instead of `*uint` and add `gorm:"not null"`

	// Additional Metadata (optional)
	// Use JSONB for flexible metadata storage in PostgreSQL
	Metadata *datatypes.JSON `gorm:"type:jsonb"` // Example: Store additional provider details
	// If you don't need Metadata right away, you can omit this field or use:
	// Metadata string `gorm:"type:text"` // Simpler text storage if JSONB isn't needed

	// Optional: Timestamps specific to payment lifecycle
	CompletedAt *time.Time // When the payment transitioned to COMPLETED
	FailedAt    *time.Time // When the payment transitioned to FAILED
}
