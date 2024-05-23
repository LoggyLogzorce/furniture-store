package entity

import (
	"furniture_store/db"
)

type ItemsOrder struct {
	Id       uint32 `json:"id" gorm:"primary_key"`
	OrderID  uint32 `json:"order_id"`
	Product  uint32 `json:"product"`
	Quantity uint8  `json:"quantity"`
}

type ItemOrderProduct struct {
	Id          uint32 `json:"id"`
	OrderID     uint   `json:"order_id"`
	ProductId   uint32 `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    uint8  `json:"quantity"`
	Price       string `json:"price"`
}

func (_ ItemsOrder) TableName() string {
	return "items_order"
}

func MigrateItemsOrder() {
	err := db.DB().AutoMigrate(ItemsOrder{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateItemsOrder)
}
