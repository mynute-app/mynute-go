package services

import (
	"fmt"

	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) UpdateOneBy(param string, value string, v interface{}) (error) {
	// Start with the base query
	query := p.DB

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch and update the first record.
	return query.First(cond, value).Updates(v).Error
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

func (p *Postgres) GetOneBy(param string, value string, v interface{}, preloads []string) (error) {
	// Start with the base query
	query := p.DB

	// Iterate over the preloads and apply each one
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch the first record by the specified parameter after applying all preloads
	return query.First(v, cond, value).Error
}

func (p *Postgres) DeleteOneBy(param string, value string, v interface{}) (error) {
	// Start with the base query
	query := p.DB

	cond := fmt.Sprintf("%s = ?", param)

	// Fetch and delete the first record.
	return query.First(cond, value).Delete(v).Error
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