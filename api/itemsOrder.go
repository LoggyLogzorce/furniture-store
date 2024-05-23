package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"strconv"
)

func UpdateItemsOrder(updatedData map[string]string) {
	uintId, err := strconv.ParseUint(updatedData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	uintProductId, err := strconv.ParseUint(updatedData["Номер товара"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	uintQuantity, err := strconv.ParseUint(updatedData["Количество"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	Id := uint32(uintId)
	ProductId := uint32(uintProductId)
	Quantity := uint8(uintQuantity)

	itemsOrder := entity.ItemsOrder{
		Product:  ProductId,
		Quantity: Quantity,
	}

	db.DB().Set("gorm:query_option", "FOR UPDATE").Where("id = ?", Id).Updates(&itemsOrder)
}

func DeleteItemsOrder(rowData map[string]string) {
	var itemsOrder entity.ItemsOrder
	uintId, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Id := uint32(uintId)
	db.DB().Where("id = ?", Id).Delete(&itemsOrder)
}
