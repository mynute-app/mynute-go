package model

type AdminRole struct {
	BaseModel
	Name        string `gorm:"type:varchar(20);uniqueIndex:idx_admin_role_name;not null" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}
