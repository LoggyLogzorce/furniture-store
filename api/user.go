package api

import (
	"furniture_store/db"
	"furniture_store/entity"
	"net/http"
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
			Expired: time.Now().Add(2 * time.Minute),
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
