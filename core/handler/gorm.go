package handler

import (
	"fmt"
	"mynute-go/core/lib"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gorm struct {
	DB            *gorm.DB
	NestedPreload *[]string
	DoNotLoad     *[]string
}

func (p *Gorm) SetNestedPreload(preloads *[]string) *Gorm {
	p.NestedPreload = preloads
	return p
}

func (p *Gorm) SetDoNotLoad(do_not_load *[]string) *Gorm {
	p.DoNotLoad = do_not_load
	return p
}

func MyGormWrapper(db *gorm.DB) *Gorm {
	return &Gorm{
		DB: db,
	}
}

// UpdateOneById updates a single record by its ID and reloads it
func (p *Gorm) UpdateOneById(value string, model any) error {
	query := p.DB

	if result := query.
		Model(model).
		Omit(clause.Associations).
		Where("id = ?", value).
		Updates(model); result.Error != nil {
		return fmt.Errorf("gorm update failed: %w", result.Error) // Wrap error
	} else if result.RowsAffected == 0 {
		var count int64
		countErr := p.DB.Model(model).Where("id = ?", value).Count(&count).Error
		if countErr != nil {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("error checking record existence: %w", countErr))
		}
		if count == 0 {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("record with id %s not found", value))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("record `id = %s` exists but was not modified maybe the changes passed are likely identical to database", value))
	}

	if p.NestedPreload != nil && len(*p.NestedPreload) > 0 {
		for _, nested_preload := range *p.NestedPreload {
			query = query.Preload(nested_preload)
		}
	}

	if err := query.First(model, "id = ?", value).Error; err != nil {
		return fmt.Errorf("gorm reload after update failed: %w", err)
	}

	return nil
}

// Create creates a new record
func (p *Gorm) Create(model any) error {
	query := p.DB

	if query.Error != nil {
		return query.Error
	}

	// Omit all associations
	query = query.Omit(clause.Associations)

	return query.Create(model).Error
}

// UpdateMany updates multiple records
func (p *Gorm) UpdateMany(v any) error {
	// Expect that `v` is a slice of records. Use GORM's Save method for bulk updates.
	// Note: GORM's `Save` updates the records if they already exist.
	return p.DB.Model(v).Updates(v).Error
}

// GetOneBy fetches a single record by a specified parameter, preloading associations.
func (p Gorm) GetOneBy(param string, value string, model any) error {
	query := p.DB // Start with the base DB instance

	// Apply Preload and update the query variable
	// The result of Preload is assigned back to query.
	query = query.Preload(clause.Associations)

	if p.NestedPreload != nil && len(*p.NestedPreload) > 0 {
		for _, nested_preload := range *p.NestedPreload {
			query = query.Preload(nested_preload)
		}
	}

	if p.DoNotLoad != nil && len(*p.DoNotLoad) > 0 {
		for _, do_not_load := range *p.DoNotLoad {
			query = query.Omit(do_not_load)
		}
	}

	// Build the WHERE condition
	cond := fmt.Sprintf("%s = ?", param)

	if err := query.First(model, cond, value).Error; err != nil {
		return err
	}

	return nil // Return nil on success
}

// ForceGetOneBy fetches a single record by a specified parameter, including soft-deleted records
func (p Gorm) ForceGetOneBy(param string, value string, model any) error {
	// Start with the base query unscoped
	query := p.DB.Unscoped()

	query = query.Preload(clause.Associations)

	if p.NestedPreload != nil && len(*p.NestedPreload) > 0 {
		for _, nested_preload := range *p.NestedPreload {
			query = query.Preload(nested_preload)
		}
	}

	if p.DoNotLoad != nil && len(*p.DoNotLoad) > 0 {
		for _, do_not_load := range *p.DoNotLoad {
			query = query.Omit(do_not_load)
		}
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch the first record by the specified parameter after applying all preloads
	return query.First(model, cond, value).Error
}

// DeleteOneById deletes a single record by its ID
func (p Gorm) DeleteOneById(value string, model any) error {
	if p.DB.Error != nil {
		return p.DB.Error
	}

	if err := p.DB.
		Model(model).
		Where("id = ?", value).
		First(model).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.DeletedError.WithError(fmt.Errorf("record with id %s not found", value))
		}
		return lib.Error.General.DeletedError.WithError(fmt.Errorf("error checking record existence: %w", err))
	}

	if err := p.DB.
		Model(model).
		Omit(clause.Associations).
		Where("id = ?", value).
		Delete(model).
		Error; err != nil {
		return lib.Error.General.DeletedError.WithError(fmt.Errorf("gorm delete failed: %w", err)) // Wrap error
	}

	return nil
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
func (p Gorm) GetAll(model any) error {
	// Verifique se model é um ponteiro antes de chamar Find()
	if model == nil {
		return fmt.Errorf("model não pode ser nil")
	}

	query := p.DB.Preload(clause.Associations)

	if p.NestedPreload != nil && len(*p.NestedPreload) > 0 {
		for _, nested_preload := range *p.NestedPreload {
			query = query.Preload(nested_preload)
		}
	}

	if query.Error != nil {
		return query.Error
	}

	// **Certifique-se de que `model` é um ponteiro ao chamar Find()**
	return query.Find(model).Error
}

// ForceGetAll fetches all records, including soft-deleted records
func (p Gorm) ForceGetAll(model any) error {
	// Start with the base query unscoped
	query := p.DB.Unscoped()

	query = query.Preload(clause.Associations)

	if p.NestedPreload != nil && len(*p.NestedPreload) > 0 {
		for _, nested_preload := range *p.NestedPreload {
			query = query.Preload(nested_preload)
		}
	}

	// Fetch all records after applying all preloads
	return query.Find(model).Error
}
