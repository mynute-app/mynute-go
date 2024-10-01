package services

import (
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) Create(v interface{}, associations []string) error {
	return p.DB.Create(v).Error
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

func (p *Postgres) Delete(v interface{}) (error) {
	return p.DB.Delete(v).Error
}

func (p *Postgres) GetOneById(v interface{}, id string, preloads []string) (error) {
		// Start with the base query
		query := p.DB

		// Iterate over the preloads and apply each one
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	
		// Fetch the first record by ID after applying all preloads
		return query.First(v, "id = ?", id).Error
}

func (p *Postgres) GetOneByName(v interface{}, name string, preloads []string) (error) {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	// Fetch the first record by name after applying all preloads
	return query.First(v, "name = ?", name).Error
}

func (p *Postgres) GetAll(v interface{}, preloads []string) (error) {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	// Fetch all records after applying all preloads
	return query.Find(v).Error
}