package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"net/http"
	"strconv"
	"time"
)

func CreateUser(user entity.User) bool {
	var checkUser entity.User
	db.DB().Where("login = ?", user.Login).First(&checkUser)
	if checkUser.Login == "" {
		db.DB().Create(&user)
		return true
	}
	return false
}

func UserRead(user entity.User) (http.Cookie, string) {
	db.DB().Where("login = ?", user.Login).First(&user)
	if user.Uid != 0 {
		token := entity.Token{
			Uid:     user.Uid,
			Token:   CreateToken(user.Login, user.Role, 2*time.Minute),
			Expired: time.Now().Add(20 * time.Minute),
		}

		cookie := http.Cookie{
			Name:  "token",
			Value: token.Token,
			Path:  "/",
		}
		return cookie, user.Role
	}

	return http.Cookie{}, ""
}

func UpdateUser(updatedData map[string]string) {
	uintVal, err := strconv.ParseUint(updatedData["uid"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Uid := uint32(uintVal)
	user := entity.User{
		Uid:      Uid,
		Name:     updatedData["Имя"],
		Login:    updatedData["Логин"],
		Password: updatedData["Пароль"],
		Role:     updatedData["Роль"],
	}
	db.DB().Save(&user)
}

func DeleteUser(rowData map[string]string) {
	var user entity.User
	uintVal, err := strconv.ParseUint(rowData["uid"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Uid := uint32(uintVal)
	db.DB().Where("uid = ?", Uid).Delete(&user)
}
