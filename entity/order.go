package entity

import (
	"furniture_store/db"
)

type Order struct {
	Id     uint32 `json:"id" gorm:"primary_key"`
	UserID uint   `json:"user_id"`
	Status string `json:"status"`
}

type OrderItem struct {
	Id          string `json:"id"`
	ProductName string `json:"product_name"`
	Quantity    uint8  `json:"quantity"`
	Price       string `json:"price"`
}

func (_ Order) TableName() string {
	return "order"
}

func MigrateOrder() {
	err := db.DB().AutoMigrate(Order{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateOrder)
}
