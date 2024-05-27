package entity

import (
	"furniture_store/db"
	"time"
)

type Order struct {
	Id         uint32    `json:"id" gorm:"primary_key"`
	UserID     uint      `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Time       time.Time `json:"time"`
}

type OrderItem struct {
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
