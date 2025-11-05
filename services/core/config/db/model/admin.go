package model

import (
	"fmt"
	"mynute-go/services/core/api/lib"
	mJSON "mynute-go/services/core/config/db/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin represents a system-wide administrator who can access all tenants
// Authentication is handled by auth service - UserID links to auth.users.id
type Admin struct {
	UserID   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"` // Primary key, references auth.users.id
	Name     string         `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Surname  string         `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"surname"`
	IsActive bool           `gorm:"default:true" json:"is_active"`
	Roles    []RoleAdmin    `gorm:"many2many:admin_role_admins;constraint:OnDelete:CASCADE;" json:"roles"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"` // Business-specific metadata (design, etc.)
}

type AdminList struct {
	Admins []Admin `json:"admins"`
	Total  int     `json:"total" example:"1"`
}

func (Admin) TableName() string  { return "public.admins" }
func (Admin) SchemaType() string { return "public" }

// BeforeCreate hook to validate before creating admin
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.UserID == uuid.Nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("user_id is required"))
	}
	if err := lib.MyCustomStructValidator(a); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate hook to validate before updating admin
func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	return lib.MyCustomStructValidator(a)
}

// GetFullAdmin loads the admin with all associations
func (a *Admin) GetFullAdmin(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}
	if err := tx.Preload("Roles").First(a, "user_id = ?", a.UserID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// HasRole checks if the admin has a specific role
func (a *Admin) HasRole(roleName string) bool {
	for _, role := range a.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsSuperAdmin checks if the admin has the superadmin role
func (a *Admin) IsSuperAdmin() bool {
	return a.HasRole("superadmin")
}

// RoleAdmin represents a role for system administrators
type RoleAdmin struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex" validate:"required,min=3,max=100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Admins      []Admin        `gorm:"many2many:admin_role_admins;foreignKey:ID;joinForeignKey:RoleAdminID;References:UserID;joinReferences:AdminUserID;constraint:OnDelete:CASCADE;" json:"admins,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RoleAdmin) TableName() string  { return "public.role_admins" }
func (RoleAdmin) SchemaType() string { return "public" }

// BeforeCreate hook to validate before creating role
func (r *RoleAdmin) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if err := lib.MyCustomStructValidator(r); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate hook to validate before updating role
func (r *RoleAdmin) BeforeUpdate(tx *gorm.DB) error {
	// Prevent changing ID
	if r.ID != uuid.Nil {
		var existingRole RoleAdmin
		if err := tx.Unscoped().Select("id").Where("id = ?", r.ID).Take(&existingRole).Error; err != nil {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("role not found"))
		}
	}
	return nil
}

// Default admin roles
var (
	RoleAdminSuperAdmin = &RoleAdmin{
		Name:        "superadmin",
		Description: "Full system access with all privileges across all tenants",
	}
	RoleAdminSupport = &RoleAdmin{
		Name:        "support",
		Description: "Customer support role with read access to tenant data",
	}
	RoleAdminAuditor = &RoleAdmin{
		Name:        "auditor",
		Description: "Audit role with read-only access to system logs and reports",
	}
)

var RoleAdmins = []*RoleAdmin{
	RoleAdminSuperAdmin,
	RoleAdminSupport,
	RoleAdminAuditor,
}

func LoadAdminRoleIDs(db *gorm.DB) error {
	for _, r := range RoleAdmins {
		var existing RoleAdmin
		if err := db.Where("name = ?", r.Name).First(&existing).Error; err != nil {
			return err
		}
		r.ID = existing.ID
	}
	return nil
}
