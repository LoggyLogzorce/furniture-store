package entity

import "furniture_store/db"

type Review struct {
	Id        uint32  `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Rating    float64 `json:"rating"`
	Comment   string  `json:"comment"`
}

func (_ Review) TableName() string {
	return "review"
}

func MigrateReview() {
	err := db.DB().AutoMigrate(Review{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateReview)
}
