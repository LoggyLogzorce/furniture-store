package entity

import "time"

type Token struct {
	Tid     uint32
	Uid     uint32
	Token   string
	Expired time.Time
}
