package DTO

import "github.com/google/uuid"

type CreateService struct {
	CompanyID   uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name        string    `json:"name" example:"Premium Consultation"`
	Description string    `json:"description" example:"A 60-minute in-depth business consultation"`
	Price       int32     `json:"price" example:"150"`
	Duration    uint      `json:"duration" example:"60"`
}

// @description	Service Full DTO
// @name			ServiceFullDTO
// @tag.name		service.full.dto
type Service struct {
	ServiceBase
	Branches    []*BranchBase   `json:"branches"`
	Employee    []*EmployeeBase  `json:"employees"`
}

// @description	Service Base DTO
// @name			ServiceBaseDTO
// @tag.name		service.base.dto
type ServiceBase struct {
	ID          uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" example:"1"`
	Name        string    `json:"name" example:"Premium Consultation"`
	Description string    `json:"description" example:"A 60-minute in-depth business consultation"`
	Price       int32     `json:"price" example:"150"`
	Duration    uint      `json:"duration" example:"60"`
}
