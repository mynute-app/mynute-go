package DTO

import (
	dJSON "mynute-go/core/src/config/api/dto/json"

	"github.com/google/uuid"
)

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
	Branches []*BranchBase   `json:"branches"`
	Employee []*EmployeeBase `json:"employees"`
}

// @description	Service Base DTO
// @name			ServiceBaseDTO
// @tag.name		service.base.dto
type ServiceBase struct {
	ID          uuid.UUID    `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID   uuid.UUID    `json:"company_id" gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" example:"1"`
	Name        string       `json:"name" example:"Premium Consultation"`
	Description string       `json:"description" example:"A 60-minute in-depth business consultation"`
	Price       int32        `json:"price" example:"150"`
	Duration    uint         `json:"duration" example:"60"`
	Design      dJSON.Design `json:"design"`
}

type ServiceID struct {
	ID uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
}

type AvailableEmployeeInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	TimeZone string    `json:"time_zone"`
}

type AvailableTime struct {
	Time        string      `json:"time"`
	EmployeesID []uuid.UUID `json:"employees"`
}

type AvailableDate struct {
	Date           string          `json:"date"`
	BranchID       uuid.UUID       `json:"branch_id"`
	AvailableTimes []AvailableTime `json:"time_slots"`
}

type ServiceAvailability struct {
	ServiceID      uuid.UUID       `json:"service_id"`
	AvailableDates []AvailableDate `json:"available_dates"`
	EmployeeInfo   []EmployeeBase  `json:"employee_info"`
	BranchInfo     []BranchBase    `json:"branch_info"`
}
