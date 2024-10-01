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

func (p *Postgres) Update(v interface{}) (error) {
	return p.DB.Save(v).Error
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