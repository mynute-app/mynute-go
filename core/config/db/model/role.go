package model

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"log"

	"gorm.io/gorm"
)

var AllowSystemRoleCreation = false

type Role struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(20);not null"`
	Description string   `gorm:"type:varchar(255)"`
	CompanyID   *uint    `gorm:"index" json:"company_id"`
	Company     Company  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE"`
	Routes      []*Route `gorm:"many2many:role_routes;constraint:OnDelete:CASCADE"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if err := r.isRoleNameReserved(tx); err != nil {
		return err
	}
	if !AllowSystemRoleCreation && r.CompanyID == nil {
		return lib.Error.Company.NotSame
	}
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	if err := r.isRoleNameReserved(tx); err != nil {
		return err
	}
	if !AllowSystemRoleCreation && r.CompanyID == nil {
		return lib.Error.Company.NotSame
	}
	return nil
}

func (r *Role) isRoleNameReserved(tx *gorm.DB) error {
	if r.CompanyID != nil {
		var count int64
		if err := tx.
			Model(&Role{}).
			Where("name = ? AND company_id IS NULL", r.Name).
			Count(&count).Error; err != nil {
			return err
		}
		return lib.Error.Role.NameReserved
	}
	return nil
}

func SeedRoles(db *gorm.DB) error {
	AllowSystemRoleCreation = true
	defer func() {
		AllowSystemRoleCreation = false
	}()

	systemRoles := []Role{
		{
			Name:        namespace.Role.Owner,
			Description: "Company Owner. Can access anything within the company's scope.",
		},
		{
			Name:        namespace.Role.GeneralManager,
			Description: "Company General Manager. Can access anything within the company's scope besides editing the company name and taxID; and deleting the company.",
		},
		{
			Name:        namespace.Role.BranchManager,
			Description: "Company Branch Manager. Can access anything within the branch's scope besides deleting, renaming and changing its address; Can also manage appointments in the branch.",
		},
		{
			Name:        namespace.Role.BranchSupervisor,
			Description: "Company Branch Supervisor. Can see anything within the branch's scope but can't change or delete anything related to branch services, employees and properties; Can also manage appointments in the branch.",
		},
		{
			Name:        namespace.Role.Employee,
			Description: "Company Employee. Can only see branches, services and appointments assigned. Besides also being able to edit its own properties.",
		},
	}

	for _, role := range systemRoles {
		var existing Role
		err := db.Where("name = ? AND company_id IS NULL", role.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	log.Println("System roles seeded successfully!")
	return nil
}

