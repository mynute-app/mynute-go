package handlers

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Gorm struct {
	DB *gorm.DB
}

func (p *Gorm) UpdateOneBy(param string, value string, model interface{}, changes interface{}, associations []string) error {
	// Start with the base query
	query := p.DB.Model(model)

	if query.Error != nil {
		return query.Error
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch the existing record
	if err := query.Where(cond, value).Error; err != nil {
		return err
	}

	// Apply the changes
	if err := query.Updates(changes).Error; err != nil {
		return err
	}

	// Get the updated record and load it into the model
	return p.GetOneBy(param, value, model, associations)
}

func (p *Gorm) Create(model interface{}) error {
	log.Printf("GORM: %v", p.DB)
	return p.DB.Create(model).Error
}

// UpdateMany updates multiple records
func (p *Gorm) UpdateMany(v interface{}) error {
	// Expect that `v` is a slice of records. Use GORM's Save method for bulk updates.
	// Note: GORM's `Save` updates the records if they already exist.
	return p.DB.Model(v).Updates(v).Error
}

func (p Gorm) GetOneBy(param string, value string, model interface{}, associations []string) error {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch the first record by the specified parameter after applying all preloads
	return query.First(model, cond, value).Error
}

func (p Gorm) DeleteOneBy(param string, value string, model interface{}) error {
	// Start with the base query
	// query := p.DB.Model(model)

	// if query.Error != nil {
	// 	return query.Error
	// }

	// cond := fmt.Sprintf("%s = ?", param)

	// // Fetch and delete the first record.
	// return query.First(cond, value).Delete(model).Error

	err := p.GetOneBy(param, value, model, nil); if err != nil {
		return err
	}

	return p.DB.Delete(model).Error
}

func (p Gorm) GetAll(model interface{}, associations []string) error {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	// Fetch all records after applying all preloads
	return query.Find(model).Error
}
