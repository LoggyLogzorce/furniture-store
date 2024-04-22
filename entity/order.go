package entity

import (
	"magazine/db"
	"time"
)

type Order struct {
	Id         uint32    `json:"id" gorm:"primary_key"`
	UserID     uint      `json:"user_id"`
	Products   uint      `json:"products"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Time       time.Time `json:"time"`
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
