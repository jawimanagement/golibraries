package models

import "gorm.io/gorm"

// model
type ProvinceModels struct {
	ProvId     int        `gorm:"<-:false;column:prov_id;primary_key" json:"prov_id"`
	ProvName   NullString `gorm:"column:prov_name" json:"prov_name"`
	LocationId int        `gorm:"column:locationid" json:"locationid"`
	Status     int        `gorm:"column:status" json:"status"`
}

type ProvinceModelsPreload struct {
	ProvId     int        `gorm:"<-:false;column:prov_id;primary_key" json:"prov_id"`
	ProvName   NullString `gorm:"column:prov_name" json:"prov_name"`
	LocationId int        `gorm:"column:locationid" json:"locationid"`
	Status     int        `gorm:"column:status" json:"status"`
}

func (c *ProvinceModels) TableName() string {
	return "cities"
}

func (c *ProvinceModelsPreload) TableName() string {
	return "cities"
}

func (c *ProvinceModels) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (c *ProvinceModels) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

func (c *ProvinceModels) BeforeDelete(tx *gorm.DB) (err error) {
	return
}
