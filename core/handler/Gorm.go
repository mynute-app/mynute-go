package handler

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gorm struct {
	DB *gorm.DB
}

// UpdateOneById updates a single record by its ID and reloads it
func (p *Gorm) UpdateOneById(value string, model any, changes any) error {
	// === Phase 1: Update ===

	// Define the target for the update operation.
	// Note: We don't need .Omit(clause.Associations) if 'changes' only contains
	//       direct column names/values of the 'model' table.
	//       If 'changes' might contain struct values that GORM could mistake
	//       for associations you want to skip updating, keep Omit.
	//       Alternatively, use .Select() in Updates for explicit control.
	// Example using Select (if 'changes' is a map[string]interface{}):
	// updateFields := make([]string, 0, len(changes))
	// for k := range changes {
	//  updateFields = append(updateFields, k)
	// }
	// result := p.DB.Model(model).Where("id = ?", value).Select(updateFields).Updates(changes)

	// Original approach (simpler if 'changes' is well-behaved):
	result := p.DB.
		Model(model).
		Where("id = ?", value).
		Omit(clause.Associations).
		Updates(changes)

	if result.Error != nil {
		return fmt.Errorf("gorm update failed: %w", result.Error) // Wrap error
	}

	// Optional but Recommended: Check if the record was actually found and updated.
	// Updates doesn't return ErrRecordNotFound if the WHERE condition matches 0 rows.
	if result.RowsAffected == 0 {
		// To be sure it wasn't found vs. data was identical, check existence separately
		var count int64
		countErr := p.DB.Model(model).Where("id = ?", value).Count(&count).Error
		if countErr != nil {
			// Error during count check, return this error
			return fmt.Errorf("failed to verify existence after 0 rows affected: %w", countErr)
		}
		if count == 0 {
			// Record definitively doesn't exist
			return gorm.ErrRecordNotFound // Return standard GORM error
		}
		// If count > 0, it means the record exists, but the data in 'changes'
		// was identical to the existing data, so RowsAffected was 0. We can continue.
		fmt.Printf("Record %s existed but was not modified (data likely identical)\n", value) // Or log warning
	}

	// === Phase 2: Reload ===

	// Perform a *new*, clean query to fetch the updated record with associations
	// Use First or Take. First is generally preferred when expecting exactly one result.
	err := p.DB.Preload(clause.Associations).First(model, "id = ?", value).Error
	if err != nil {
		// If it fails here (e.g., record not found), it reflects the state *after* the update attempt.
		return fmt.Errorf("gorm reload after update failed: %w", err) // Wrap error
	}

	// 'model' pointer now holds the latest data with preloaded associations
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
func (p Gorm) GetAll(model any) error {
	// Verifique se model é um ponteiro antes de chamar Find()
	if model == nil {
		return fmt.Errorf("model não pode ser nil")
	}

	query := p.DB.Preload(clause.Associations)

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

	// Fetch all records after applying all preloads
	return query.Find(model).Error
}
