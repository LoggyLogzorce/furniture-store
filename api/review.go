package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"strconv"
)

func DeleteReview(rowData map[string]string) {
	var review entity.Review
	uintId, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Id := uint32(uintId)
	db.DB().Where("id = ?", Id).Delete(&review)
}
