package models

import "gorm.io/gorm"

type CityModels struct {
	CityID   int    `gorm:"column:city_id;primary_key" json:"city_id"`
	CityName string `gorm:"column:city_name" json:"city_name"`
	ProvID   int    `gorm:"column:prov_id" json:"prov_id"`
}

type CityModelPreload struct {
	CityID       int            `gorm:"column:city_id;primary_key" json:"city_id"`
	CityName     string         `gorm:"column:city_name" json:"city_name"`
	ProvID       int            `gorm:"column:prov_id" json:"prov_id"`
	ProvinceInfo ProvinceModels `gorm:"->;foreignKey:ID;references:ProvID" json:"province_info"`
}

func (p *CityModels) TableName() string {
	return "cities"
}

func (p *CityModelPreload) TableName() string {
	return "cities"
}

func (p *CityModels) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (p *CityModels) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}
func (p *CityModels) AfterUpdate(tx *gorm.DB) (err error) {
	return
}
func (p *CityModels) BeforeDelete(tx *gorm.DB) (err error) {
	return
}
