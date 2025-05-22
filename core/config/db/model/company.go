package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Uploader"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Company struct {
	BaseModel
	LegalName  string             `gorm:"not null;unique" json:"legal_name"` // Razão Social
	TradeName  string             `gorm:"not null" json:"trade_name"`        // Nome Fantasia
	TaxID      string             `gorm:"not null;uniqueIndex" json:"tax_id"`
	SchemaName string             `gorm:"type:varchar(100);not null;uniqueIndex" json:"schema_name"`
	Subdomains []*Subdomain       `gorm:"constraint:OnDelete:CASCADE;" json:"subdomains"`                        // One-to-many relationship with Subdomain
	Sectors    []*Sector          `gorm:"many2many:company_sectors;constraint:OnDelete:CASCADE;" json:"sectors"` // Many-to-many relationship with Sector
	Design     mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

func (Company) TableName() string { return "public.companies" }
func (Company) SchemaType() string { return "public" }

func (c *Company) GenerateSchemaName() string {
	return "company" + "_" + c.ID.String()
}

func (c *Company) BeforeUpdate(tx *gorm.DB) error {
	return c.CheckDesignImageOwnershipBeforeUpdate(tx)
}

func (c *Company) CheckDesignImageOwnershipBeforeUpdate(tx *gorm.DB) error {
	// carregar a versão antiga no banco
	var old Company
	if err := tx.Unscoped().Model(&Company{}).Where("id = ?", c.ID).First(&old).Error; err != nil {
		return err
	}

	oldImgs := old.Design.Images
	newImgs := c.Design.Images

	// map[string]string{nomeCampo: url}
	imagePairs := map[string][2]string{
		"logo":      {oldImgs.LogoURL, newImgs.LogoURL},
		"banner":    {oldImgs.BannerURL, newImgs.BannerURL},
		"favicon":   {oldImgs.FaviconURL, newImgs.FaviconURL},
		"background": {oldImgs.BackgroundURL, newImgs.BackgroundURL},
	}

	for key, pair := range imagePairs {
		oldURL := pair[0]
		newURL := pair[1]

		if oldURL != newURL && newURL != "" {
			if !imageBelongsToCompany(c.ID, newURL) {
				return lib.Error.General.BadRequest.WithError(
					fmt.Errorf("image in field '%s' does not belong to this company", key),
				)
			}
		}
	}

	return nil
}

func imageBelongsToCompany(companyID uuid.UUID, imageURL string) bool {
	filename := myUploader.ExtractFilenameFromURL(imageURL)
	return strings.Contains(filename, companyID.String())
}



func (c *Company) MigrateSchema(tx *gorm.DB) error {
	c.SchemaName = c.GenerateSchemaName()

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.Save(c).Error; err != nil {
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

	// Migrate tenant specific tables
	for _, model := range TenantModels {
		if err := tx.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
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

	if err := c.MigrateSchema(tx); err != nil {
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