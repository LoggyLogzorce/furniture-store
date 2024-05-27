package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"strconv"
)

func UpdateProduct(updatedData map[string]string) {
	uintId, err := strconv.ParseUint(updatedData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	uintCategory, err := strconv.ParseUint(updatedData["Категория"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	Id := uint32(uintId)
	Category := uint32(uintCategory)

	product := entity.Product{
		Id:          Id,
		Category:    Category,
		Name:        updatedData["Название"],
		Price:       updatedData["Цена"],
		Description: updatedData["Описание"],
		Image:       updatedData["Изображение"],
	}
	db.DB().Save(&product)
}

func DeleteProduct(rowData map[string]string) {
	var product entity.Product
	uintId, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Id := uint32(uintId)
	db.DB().Where("id = ?", Id).Delete(&product)
}

func AddProduct(rowData map[string]string) {
	uintCategory, err := strconv.ParseUint(rowData["category"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	var product entity.Product
	product.Category = uint32(uintCategory)
	product.Name = rowData["name"]
	product.Price = rowData["price"]
	product.Description = rowData["description"]
	product.Image = rowData["image"]
	db.DB().Create(&product)
}
