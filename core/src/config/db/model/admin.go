package model

import (
	"fmt"
	mJSON "mynute-go/core/src/config/db/model/json"
	"mynute-go/core/src/lib"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin represents a system-wide administrator who can access all tenants
type Admin struct {
	BaseModel
	Name     string         `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Email    string         `gorm:"type:varchar(100);uniqueIndex" validate:"required,email" json:"email"`
	Password string         `gorm:"type:varchar(255)" validate:"required,myPasswordValidation" json:"password"`
	IsActive bool           `gorm:"default:true" json:"is_active"`
	Roles    []RoleAdmin    `gorm:"many2many:admin_role_admins;constraint:OnDelete:CASCADE;" json:"roles"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
}

func (Admin) TableName() string  { return "public.admins" }
func (Admin) SchemaType() string { return "public" }

// BeforeCreate hook to validate and hash password before creating admin
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(a); err != nil {
		return err
	}
	if err := a.HashPassword(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate hook to validate and hash password before updating admin
func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	if a.Password != "" {
		var dbAdmin Admin
		tx.First(&dbAdmin, "id = ?", a.ID)
		// Skip hashing if password hasn't changed
		if a.Password == dbAdmin.Password || a.MatchPassword(dbAdmin.Password) {
			return nil
		}
		if err := lib.ValidatorV10.Var(a.Password, "myPasswordValidation"); err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("password invalid"))
			} else {
				return lib.Error.General.InternalError.WithError(err)
			}
		}
		return a.HashPassword()
	}
	return nil
}

// MatchPassword compares the provided plain text password with the hashed password
func (a *Admin) MatchPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plainPassword))
	return err == nil
}

// HashPassword generates a bcrypt hash of the password
func (a *Admin) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hash)
	return nil
}

// GetFullAdmin loads the admin with all associations
func (a *Admin) GetFullAdmin(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}
	aID := a.ID.String()
	if err := tx.Preload("Roles").First(a, "id = ?", aID).Error; err != nil {
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
	BaseModel
	Name        string  `gorm:"type:varchar(100);uniqueIndex" validate:"required,min=3,max=100" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Admins      []Admin `gorm:"many2many:admin_role_admins;constraint:OnDelete:CASCADE;" json:"admins,omitempty"`
}

func (RoleAdmin) TableName() string  { return "public.role_admins" }
func (RoleAdmin) SchemaType() string { return "public" }

// BeforeCreate hook to validate before creating role
func (r *RoleAdmin) BeforeCreate(tx *gorm.DB) error {
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
