package api

import (
	"fmt"
	"furniture_store/db"
	"furniture_store/entity"
	"log"
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

func DeleteItemsOrder(rowData map[string]string, token string) {
	fmt.Println(rowData)
	Id, err := strconv.ParseUint(rowData["id"], 10, 32)
	if err != nil {
		log.Println("Ошибка конвертации:", err)
		return
	}

	product := uint32(Id)
	login, err := GetLoginFromToken(token)
	fmt.Println(login)
	if err != nil {
		log.Println(err)
	}
	var user entity.User
	var order entity.Order
	var itemsOrder entity.ItemsOrder
	itemsOrder.Product = product

	db.DB().Where("login = ?", login).Find(&user)
	db.DB().Where("user_id = ?", user.Uid).Find(&order)

	fmt.Println(user)
	db.DB().Where("order_id = ? and product = ?", order.Id, product).Delete(&itemsOrder)
}

func AddCart(rowData map[string]string, token string) {
	uintId, err := strconv.ParseUint(rowData["productId"], 10, 32)
	if err != nil {
		log.Println("Ошибка конвертации:", err)
		return
	}

	Id := uint32(uintId)
	login, err := GetLoginFromToken(token)
	if err != nil {
		log.Println(err)
	}
	var user entity.User
	var order entity.Order
	var itemOrder entity.ItemsOrder
	itemOrder.Product = Id

	db.DB().Where("login = ?", login).Find(&user)
	db.DB().Where("user_id = ?", user.Uid).Find(&order)
	if order.Id == 0 {
		order.UserID = uint(user.Uid)
		order.Status = "Корзина"
		db.DB().Create(&order)
	}

	db.DB().Where("order_id = ? and product = ?", order.Id, itemOrder.Product).Find(&itemOrder)
	if itemOrder.Id == 0 {
		itemOrder.OrderID = order.Id
		itemOrder.Quantity += 1
		db.DB().Create(&itemOrder)
		return
	}
	itemOrder.Quantity += 1
	db.DB().Where("order_id = ? and product = ?", order.Id, itemOrder.Product).Save(&itemOrder)
	return
}
