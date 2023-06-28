package models

import "gorm.io/gorm"

type CityModels struct {
	CityId   int        `gorm:"<-:false;column:city_id;primary_key" json:"city_id"`
	CityName NullString `gorm:"column:city_name" json:"city_name"`
	ProvId   int        `gorm:"column:prov_id" json:"prov_id"`
}

// city model preload
type CityModelsPreload struct {
	CityId   int                   `gorm:"column:city_id;primary_key" json:"city_id"`
	CityName NullString            `gorm:"column:city_name" json:"city_name"`
	ProvId   int                   `gorm:"column:prov_id" json:"prov_id"`
	ProvInfo ProvinceModelsPreload `gorm:"foreignKey:ProvId;references:ProvId" json:"province_info"`
}

func (c *CityModels) TableName() string {
	return "cities"
}

func (c *CityModelsPreload) TableName() string {
	return "cities"
}

func (c *CityModels) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (c *CityModels) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

func (c *CityModels) BeforeDelete(tx *gorm.DB) (err error) {
	return
}
