package entity

import "magazine/db"

type User struct {
	Uid      uint32 `json:"uid" gorm:"primaryKey"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Addm     bool   `json:"addm"`
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
