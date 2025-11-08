package model

type Role struct {
	BaseModel             // Adds ID (uint), CreatedAt, UpdatedAt, DeletedAt
	Name        string    `gorm:"type:varchar(20);not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}

