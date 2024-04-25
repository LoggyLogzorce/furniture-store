package entity

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	Tid     uint32
	Uid     uint32
	Token   string
	Expired time.Time
}

type Claims struct {
	Login string `json:"login"`
	Role  bool   `json:"role"`
	jwt.StandardClaims
}
