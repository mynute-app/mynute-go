package model

import (
	"github.com/google/uuid"
)

type Subdomain struct {
	BaseModel
	Name      string    `gorm:"type:varchar(36);not null;uniqueIndex" json:"name"`
	CompanyID uuid.UUID `gorm:"not null;index" json:"company_id"`
	Company   *Company  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`
}

func (Subdomain) TableName() string {
	return "public.subdomains"
}