package model

import "gorm.io/gorm"

type Property struct {
	gorm.Model
	Name         string   `json:"name" gorm:"unique;not null"`
	Description  string   `json:"description"`
	ResourceName string   `json:"resource_name"`
	Resource     Resource `gorm:"foreignKey:ResourceName;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

var CompanyCNPJ = &Property{
	Name:         "company_cnpj",
	Description:  "CNPJ of the company",
	ResourceName: CompanyResource.Name,
}
