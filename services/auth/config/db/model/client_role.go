package model

type ClientRole struct {
	BaseModel
	Name        string `gorm:"type:varchar(20);uniqueIndex:idx_client_role_name;not null" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}
