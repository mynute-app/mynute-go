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
	IsSystemRole bool       `gorm:"not null;default:false" json:"is_system_role"`
}

func (Role) TableName() string { return "public.roles" }
func (Role) SchemaType() string { return "public" }

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if err := r.isRoleNameReserved(tx); err != nil {
		return err
	}
	if r.ID == uuid.Nil && !AllowSystemRoleCreation && r.CompanyID == nil {
		return lib.Error.Role.NilCompanyID
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
			Where("name = ?", r.Name).
			Count(&count).Error; err != nil {
			return err
		}
		return lib.Error.Role.NameReserved
	}
	return nil
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

// func SeedRoles(db *gorm.DB) ([]*Role, error) {
// 	AllowSystemRoleCreation = true
// 	tx := db.Begin()
// 	seededCount := 0
// 	updatedCount := 0
// 	defer func() {
// 		AllowSystemRoleCreation = false
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			log.Printf("Panic occurred during policy seeding: %v", r)
// 		}
// 		if err := tx.Commit().Error; err != nil {
// 			log.Printf("Failed to commit transaction: %v", err)
// 		}
// 		log.Printf("System Roles seeded successfully. New: %d, Existing/Updated: %d", seededCount, updatedCount)
// 	}()
// 	for _, newRole := range Roles {
// 		if newRole == nil {
// 			panic("Encountered nil role in Roles slice")
// 		}
// 		var oldRole Role
// 		err := tx.Where("name = ?", newRole.Name).First(oldRole).Error
// 		if err == gorm.ErrRecordNotFound {
// 			if err := tx.Create(newRole).Error; err != nil {
// 				return nil, fmt.Errorf("failed to create role %s: %w", newRole.Name, err)
// 			}
// 			seededCount++
// 		} else if err != nil {
// 			return nil, fmt.Errorf("failed to check if role %s exists: %w", newRole.Name, err)
// 		} else {
// 			if err := tx.Model(newRole).Where("id = ?", ).Updates(newRole).Error; err != nil {
// 				return nil, fmt.Errorf("failed to update role %s: %w", newRole.Name, err)
// 			}
// 			updatedCount++
// 		}
// 	}
// 	return Roles, nil
// }
