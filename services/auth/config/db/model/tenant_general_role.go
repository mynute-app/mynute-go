package model

type TenantGeneralRole struct {
	BaseModel
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}