package entity

import (
	"furniture_store/db"
)

type Category struct {
	Id   uint32 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (_ Category) TableName() string {
	return "category"
}

func MigrateCategory() {
	err := db.DB().AutoMigrate(Category{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateCategory)
}
