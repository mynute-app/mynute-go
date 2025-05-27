package model

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

var AllowSystemRoleCreation = false

type Role struct {
	BaseModel               // Adds ID (uint), CreatedAt, UpdatedAt, DeletedAt
	Name         string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_role_name_company,priority:1" json:"name"`
	Description  string     `json:"description"`
	CompanyID    *uuid.UUID `gorm:"index;uniqueIndex:idx_role_name_company,priority:2" json:"company_id"` // Null for system roles
	Company      *Company   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`     // BelongsTo Company
}

func (Role) TableName() string  { return "public.roles" }
func (Role) SchemaType() string { return "public" }

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if is, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if is {
		return lib.Error.Role.NameReserved
	}
	if r.ID == uuid.Nil && !AllowSystemRoleCreation && r.CompanyID == nil {
		return lib.Error.Role.NilCompanyID
	}
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	if nameIsReserved, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if nameIsReserved && tx.Statement.Changed("name") {
		return lib.Error.Role.NameReserved
	}
	return nil
}

func (r *Role) isRoleNameReserved(tx *gorm.DB) (bool, error) {
	var count int64
	if err := tx.
		Model(&Role{}).
		Where("name = ? AND company_id IS NULL", r.Name).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

var SystemRoleOwner = &Role{
	Name:        namespace.Role.Owner,
	Description: "Company Owner. Can access anything within the company's scope.",
}

var SystemRoleGeneralManager = &Role{
	Name:        namespace.Role.GeneralManager,
	Description: "Company General Manager. Can access anything within the company's scope besides editing the company name and taxID; and deleting the company.",
}

var SystemRoleBranchManager = &Role{
	Name:        namespace.Role.BranchManager,
	Description: "Company Branch Manager. Can access anything within the branch's scope besides deleting, renaming and changing its address; Can also manage appointments in the branch.",
}

var SystemRoleBranchSupervisor = &Role{
	Name:        namespace.Role.BranchSupervisor,
	Description: "Company Branch Supervisor. Can see anything within the branch's scope but can't change or delete anything related to branch services, employees and properties; Can also manage appointments in the branch.",
}

// --- Combine all system roles into a slice for seeding ---
var Roles = []*Role{
	SystemRoleOwner,
	SystemRoleGeneralManager,
	SystemRoleBranchManager,
	SystemRoleBranchSupervisor,
}

func LoadSystemRoleIDs(db *gorm.DB) error {
	for _, r := range Roles {
		var existing Role
		if err := db.
			Where("name = ? AND company_id IS NULL", r.Name).
			First(&existing).Error; err != nil {
			return err
		}
		r.ID = existing.ID
	}
	return nil
}
