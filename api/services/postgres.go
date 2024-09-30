package services

import (
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) Create(v interface{}) (error) {
	result := p.DB.Create(&v)
	return result.Error
}

func (p *Postgres) Update(v interface{}) (error) {
	result := p.DB.Save(&v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Postgres) Delete(v interface{}) (error) {
	result := p.DB.Delete(&v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Postgres) GetOneById(v interface{}, id string) (error) {
	result := p.DB.First(&v, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}