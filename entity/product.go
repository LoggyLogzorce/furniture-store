package entity

import (
	"furniture_store/db"
)

type Product struct {
	Id          uint32 `json:"id" gorm:"primaryKey"`
	Category    uint   `json:"category"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

func (_ Product) TableName() string {
	return "product"
}

func MigrateProduct() {
	err := db.DB().AutoMigrate(Product{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateProduct)
}
