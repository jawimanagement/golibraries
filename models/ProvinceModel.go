package models

type ProvinceModels struct {
	ID         int        `gorm:"column:prov_id;primary_key" json:"prov_id"`
	ProvName   NullString `gorm:"column:prov_name" json:"prov_name"`
	LocationID int        `gorm:"column:locationid" json:"locationid"`
	Status     int        `gorm:"column:status" json:"status"`
}

func (*ProvinceModels) TableName() string {
	return "provinces"
}
