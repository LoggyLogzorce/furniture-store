package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"strconv"
	"time"
)

func UpdateOrder(updatedData map[string]string) {
	uintId, err := strconv.ParseUint(updatedData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	uintUid, err := strconv.ParseUint(updatedData["UserID"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	floatSum, err := strconv.ParseFloat(updatedData["Сумма"], 64)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	Id := uint32(uintId)
	Uid := uint(uintUid)
	Sum := floatSum

	order := entity.Order{
		Id:         Id,
		UserID:     Uid,
		TotalPrice: Sum,
		Status:     updatedData["Статус"],
		Time:       time.Now(),
	}
	db.DB().Save(&order)
}

func DeleteOrder(rowData map[string]string) {
	var order entity.Order
	uintId, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Id := uint32(uintId)
	db.DB().Where("id = ?", Id).Delete(&order)
}
