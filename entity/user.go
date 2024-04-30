package entity

import "furniture_store/db"

type User struct {
	Uid      uint32 `json:"uid" gorm:"primaryKey"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) TableName() string {
	return "user"
}

func MigrateUser() {
	err := db.DB().AutoMigrate(User{})
	if err != nil {
		panic(err)
	}
}

func init() {
	db.Add(MigrateUser)
}
