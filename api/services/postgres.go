package services

import (
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) Create(v interface{}) (error) {
	return p.DB.Create(&v).Error
}

func (p *Postgres) Update(v interface{}) (error) {
	return p.DB.Save(&v).Error
}

func (p *Postgres) Delete(v interface{}) (error) {
	return p.DB.Delete(&v).Error
}

func (p *Postgres) GetOneById(v interface{}, id string) (error) {
	return p.DB.First(&v, id).Error
}