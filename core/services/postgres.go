package services

import (
	"fmt"

	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

// func (p *Postgres) UpdateOneBy(param string, value string, model interface{}, changes interface{}, associations []string) error {
// 	if err := p.GetOneBy(param, value, model, associations); err != nil {
// 		return err
// 	}

// 	if err := lib.MergeMapIntoInterface(model, changes) ; err != nil {
// 		return err
// 	}

// }

func (p *Postgres) UpdateOneBy(param string, value string, model interface{}, changes interface{}, associations []string) error {
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

	// Use GORM's Updates method to update the record
	if err := query.Updates(changes).Error; err != nil {
		return err
	}

	// Get the updated record and load it into the model
	return p.GetOneBy(param, value, model, associations)
}

func (p *Postgres) Create(model interface{}, associations []string) error {
	return p.DB.Create(model).Error
}

// UpdateOne updates a single record
func (p *Postgres) UpdateOne(v interface{}) error {
	// Use GORM's Save method to update the record
	return p.DB.Save(v).Error
}

// UpdateMany updates multiple records
func (p *Postgres) UpdateMany(v interface{}) error {
	// Expect that `v` is a slice of records. Use GORM's Save method for bulk updates.
	// Note: GORM's `Save` updates the records if they already exist.
	return p.DB.Model(v).Updates(v).Error
}

func (p *Postgres) GetOneBy(param string, value string, model interface{}, associations []string) error {
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

func (p *Postgres) DeleteOneBy(param string, value string, model interface{}) error {
	// Start with the base query
	query := p.DB

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch and delete the first record.
	return query.First(cond, value).Delete(model).Error
}

func (p *Postgres) GetAll(model interface{}, associations []string) error {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range associations {
		query = query.Preload(preload)
	}

	// Fetch all records after applying all preloads
	return query.Find(model).Error
}
