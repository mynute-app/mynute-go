package model

import (
	"fmt"
	mJSON "mynute-go/core/src/config/db/model/json"
	"mynute-go/core/src/lib"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Company struct {
	BaseModel
	LegalName  string             `gorm:"type:varchar(100);uniqueIndex" validate:"required,min=3,max=100" json:"legal_name"`
	TradeName  string             `gorm:"type:varchar(100);uniqueIndex" validate:"required,min=3,max=100" json:"trade_name"`
	TaxID      string             `gorm:"type:varchar(100);uniqueIndex" validate:"required,min=3,max=100" json:"tax_id"`
	SchemaName string             `gorm:"type:varchar(100);uniqueIndex" json:"schema_name"`
	Subdomains []*Subdomain       `gorm:"constraint:OnDelete:CASCADE;" json:"subdomains"`
	Sectors    []*Sector          `gorm:"many2many:company_sectors;constraint:OnDelete:CASCADE;" json:"sectors"`
	Design     mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

func (Company) TableName() string  { return "public.companies" }
func (Company) SchemaType() string { return "public" }

func (c *Company) GenerateSchemaName() string {
	return "company" + "_" + c.ID.String()
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(c); err != nil {
		return err
	}
	return nil
}

func (c *Company) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (c *Company) AfterCreate(tx *gorm.DB) error {
	// Update schema_name using the database-generated UUID
	// Use Updates with map to avoid triggering hooks and validation
	schema_name := c.GenerateSchemaName()
	if err := tx.Model(c).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"schema_name": schema_name,
	}).Error; err != nil {
		return fmt.Errorf("failed to update schema_name: %w", err)
	}

	// Update the in-memory object
	c.SchemaName = schema_name

	if err := c.MigrateSchema(tx); err != nil {
		return err
	}
	return nil
}

func (c *Company) MigrateSchema(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	// Create the schema in the lib
	if err := c.CreateSchema(tx); err != nil {
		return err
	}

	// Set the search path to the new schema
	if err := lib.ChangeToCompanySchema(tx, c.SchemaName); err != nil {
		return err
	}

	// Migrate tenant specific tables without foreign key constraints
	// Cross-schema FK constraints cause issues, so we disable them and rely on application-level integrity
	migrator := tx.Migrator()
	for _, model := range TenantModels {
		if err := migrator.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	// Create employee_roles join table manually since we skipped Role migration
	// This table links employees in company schema to roles in public schema
	createEmployeeRolesSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s".employee_roles (
			employee_id uuid NOT NULL,
			role_id uuid NOT NULL,
			PRIMARY KEY (employee_id, role_id),
			CONSTRAINT fk_employee_roles_employee FOREIGN KEY (employee_id) REFERENCES "%s".employees(id) ON DELETE CASCADE
		)
	`, c.SchemaName, c.SchemaName)
	if err := tx.Exec(createEmployeeRolesSQL).Error; err != nil {
		return fmt.Errorf("failed to create employee_roles table: %w", err)
	}

	roleViewSQL := fmt.Sprintf(`CREATE OR REPLACE VIEW "%s".roles AS SELECT * FROM public.roles`, c.SchemaName)
	if err := tx.Exec(roleViewSQL).Error; err != nil {
		return err
	}

	clientViewSQL := fmt.Sprintf(`CREATE OR REPLACE VIEW "%s".clients AS SELECT * FROM public.clients`, c.SchemaName)
	if err := tx.Exec(clientViewSQL).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) Create(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.Create(c).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) Delete(tx *gorm.DB) error {
	if err := c.Refresh(tx); err != nil {
		return err
	}

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.Delete(c).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) CreateSchema(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if c.SchemaName == "" {
		return lib.Error.Company.SchemaIsEmpty
	}

	if err := tx.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, c.SchemaName)).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) CreateOwner(tx *gorm.DB, owner *Employee) error {
	if err := lib.ChangeToCompanySchema(tx, c.SchemaName); err != nil {
		return err
	}

	if c.ID == uuid.Nil {
		return lib.Error.Company.CouldNotCreateOwner.WithError(fmt.Errorf("company ID is empty"))
	}

	owner.CompanyID = c.ID

	if err := tx.Create(&owner).Error; err != nil {
		return err
	}

	var ownerRole Role

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.First(&ownerRole, "id = ?", SystemRoleOwner.ID).Error; err != nil {
		return err
	}

	if err := lib.ChangeToCompanySchema(tx, c.SchemaName); err != nil {
		return err
	}

	// TODO: Check why this employee_roles implementation is not working.
	if err := tx.Exec("INSERT INTO employee_roles (role_id, employee_id) VALUES (?, ?)", ownerRole.ID, owner.ID).Error; err != nil {
		return err
	}

	if err := tx.Model(&owner).Preload("Roles").Where("id = ?", owner.ID.String()).First(&owner).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) Refresh(tx *gorm.DB) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.
		Model(c).
		Preload(clause.Associations).
		Where("id = ?", c.ID.String()).
		First(c).Error; err != nil {
		return err
	}

	return nil
}

func (c *Company) GetFullCompany(tx *gorm.DB) (*CompanyMerged, error) {

	if err := c.Refresh(tx); err != nil {
		return nil, err
	}

	if err := lib.ChangeToCompanySchema(tx, c.SchemaName); err != nil {
		return nil, err
	}

	fullCompany := &CompanyMerged{
		Company: *c,
	}

	// Get Branches
	var branches []Branch
	if err := tx.Model(&Branch{}).Find(&branches).Error; err != nil {
		return nil, err
	}

	fullCompany.Branches = branches

	// Get Employees
	var employees []Employee
	if err := tx.Model(&Employee{}).Find(&employees).Error; err != nil {
		return nil, err
	}

	fullCompany.Employees = employees

	// Get Services
	var services []Service
	if err := tx.Model(&Service{}).Find(&services).Error; err != nil {
		return nil, err
	}

	fullCompany.Services = services

	return fullCompany, nil
}

func (c *Company) AddSubdomain(tx *gorm.DB, subdomain *Subdomain) error {
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	var subdomainCount int64
	if err := tx.
		Model(&Subdomain{}).
		Where("name = ?", subdomain.Name).
		Count(&subdomainCount).
		Error; err != nil {
		return err
	}

	if subdomainCount > 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("subdomain %s already exists", subdomain.Name))
	}

	if err := tx.Create(subdomain).Error; err != nil {
		return err
	}

	return nil
}

type CompanyMerged struct {
	Company
	Branches  []Branch   `json:"branches"`
	Employees []Employee `json:"employees"`
	Services  []Service  `json:"services"`
}
