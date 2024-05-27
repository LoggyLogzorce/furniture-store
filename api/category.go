package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"strconv"
)

func UpdateCategory(updatedData map[string]string) {
	uintId, err := strconv.ParseUint(updatedData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	Id := uint32(uintId)

	category := entity.Category{
		Id:   Id,
		Name: updatedData["Категория"],
	}
	db.DB().Save(&category)
}

func DeleteCategory(rowData map[string]string) {
	var category entity.Category
	uintId, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Id := uint32(uintId)
	db.DB().Where("id = ?", Id).Delete(&category)
}

func AddCategory(rowData map[string]string) {
	category := entity.Category{
		Name: rowData["category"],
	}
	db.DB().Create(&category)
}
