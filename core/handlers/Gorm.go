package handlers

import (
	"fmt"

	"gorm.io/gorm"
)

type Gorm struct {
	DB *gorm.DB
}

// UpdateOneById updates a single record by its ID
func (p *Gorm) UpdateOneById(value string, model any, changes any, associations []string) error {
	// Start with the base query
	query := p.DB.Model(model)

	if query.Error != nil {
		return query.Error
	}

	cond := fmt.Sprintf("%s = ?", "id")

	// Fetch the existing record
	if err := query.Where(cond, value).Error; err != nil {
		return err
	}

	// Apply the changes
	if err := query.Updates(changes).Error; err != nil {
		return err
	}

	// Get the updated record and load it into the model
	return p.GetOneBy("id", value, model, associations)
}

// Create creates a new record
func (p *Gorm) Create(model any) error {
	query := p.DB

	if query.Error != nil {
		return query.Error
	}

	return query.Create(model).Error
}

// UpdateMany updates multiple records
func (p *Gorm) UpdateMany(v any) error {
	// Expect that `v` is a slice of records. Use GORM's Save method for bulk updates.
	// Note: GORM's `Save` updates the records if they already exist.
	return p.DB.Model(v).Updates(v).Error
}

// GetOneBy fetches a single record by a specified parameter
// GetOneBy fetches a single record by a specified parameter
func (p Gorm) GetOneBy(param string, value string, model any, associations []string) error {
	query := p.DB

	if query.Error != nil {
		return query.Error
	}

	// Forcefully preload associations
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Run query and return result
	return query.First(model, cond, value).Error
}

// ForceGetOneBy fetches a single record by a specified parameter, including soft-deleted records
func (p Gorm) ForceGetOneBy(param string, value string, model any, associations []string) error {
	// Start with the base query unscoped
	query := p.DB.Unscoped()

	// Iterate over the preloads and apply each one
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch the first record by the specified parameter after applying all preloads
	return query.First(model, cond, value).Error
}

// DeleteOneById deletes a single record by its ID
func (p Gorm) DeleteOneById(value string, model any) error {
	query := p.DB

	if query.Error != nil {
		return query.Error
	}

	cond := fmt.Sprintf("%s = ?", "id")

	return query.Model(model).Delete(cond, value).Error
}

// ForceDeleteOneById deletes a single record by its ID, including soft-deleted records
func (p Gorm) ForceDeleteOneById(value string, model any) error {
	query := p.DB.Unscoped().Model(model)

	if query.Error != nil {
		return query.Error
	}

	cond := fmt.Sprintf("%s = ?", "id")

	return query.Delete(cond, value).Error
}

// GetAll fetches all records
func (p Gorm) GetAll(model any, associations []string) error {
	// Verifique se model é um ponteiro antes de chamar Find()
	if model == nil {
		return fmt.Errorf("model não pode ser nil")
	}

	query := p.DB

	// Aplicar os Preloads
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	// **Certifique-se de que `model` é um ponteiro ao chamar Find()**
	return query.Find(model).Error
}

// ForceGetAll fetches all records, including soft-deleted records
func (p Gorm) ForceGetAll(model any, associations []string) error {
	// Start with the base query unscoped
	query := p.DB.Unscoped()

	// Iterate over the preloads and apply each one
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	// Fetch all records after applying all preloads
	return query.Find(model).Error
}
