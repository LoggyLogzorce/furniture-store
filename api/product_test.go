package api

import (
	"furniture_store/db"
	"furniture_store/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProduct(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	// Подготовка тестовых данных
	updatedData := map[string]string{
		"id":          "1",
		"Категория":   "1",
		"Название":    "Новое название",
		"Цена":        "100",
		"Описание":    "Новое описание",
		"Изображение": "новое_изображение.jpg",
	}

	// Создание тестового продукта в базе данных
	originalProduct := entity.Product{
		Id:          1,
		Category:    1,
		Name:        "Стул",
		Price:       "50",
		Description: "Описание стула",
		Image:       "стул.jpg",
	}
	err := db.DB().Create(&originalProduct).Error
	assert.NoError(t, err)

	// Вызов тестируемой функции
	UpdateProduct(updatedData)

	// Получение обновленного продукта из базы данных
	updatedProduct, err := GetProductByID(originalProduct.Id)
	assert.NoError(t, err)
	assert.NotNil(t, updatedProduct, "Продукт должен быть найден в базе данных")

	// Проверка, что данные продукта обновлены
	assert.Equal(t, updatedData["Название"], updatedProduct.Name)
	assert.Equal(t, updatedData["Цена"], updatedProduct.Price)
	assert.Equal(t, updatedData["Описание"], updatedProduct.Description)
	assert.Equal(t, updatedData["Изображение"], updatedProduct.Image)
	db.DB().Exec("TRUNCATE TABLE product RESTART IDENTITY CASCADE")
}

func TestDeleteProduct(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	// Подготовка тестовых данных
	testRowData := map[string]string{
		"id": "1",
	}

	// Создание тестового продукта
	testProduct := entity.Product{
		Id:          1,
		Category:    1,
		Name:        "Test Product",
		Price:       "50",
		Description: "Test Description",
		Image:       "test.jpg",
	}

	err := db.DB().Create(&testProduct).Error
	assert.NoError(t, err)

	// Вызов тестируемой функции
	DeleteProduct(testRowData)

	// Проверка, что продукт удален
	deletedProduct, err := GetProductByID(testProduct.Id)
	assert.NoError(t, err)
	assert.Nil(t, deletedProduct, "Продукт должен быть удален из базы данных")
	db.DB().Exec("TRUNCATE TABLE product RESTART IDENTITY CASCADE")
}

func TestAddProduct(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	// Подготовка тестовых данных
	testRowData := map[string]string{
		"category":    "1",
		"name":        "Новый продукт",
		"price":       "100",
		"description": "Описание нового продукта",
		"image":       "новое_изображение.jpg",
	}

	// Вызов тестируемой функции
	AddProduct(testRowData)

	// Проверка, что продукт был успешно добавлен в базу данных
	addedProduct, err := GetProductByName(testRowData["name"])
	assert.NoError(t, err)
	assert.NotNil(t, addedProduct, "Продукт должен быть найден в базе данных")

	// Проверка, что данные продукта были правильно добавлены
	assert.Equal(t, uint32(1), addedProduct.Category)
	assert.Equal(t, testRowData["name"], addedProduct.Name)
	assert.Equal(t, testRowData["price"], addedProduct.Price)
	assert.Equal(t, testRowData["description"], addedProduct.Description)
	assert.Equal(t, testRowData["image"], addedProduct.Image)
	db.DB().Exec("TRUNCATE TABLE product RESTART IDENTITY CASCADE")
}

func GetProductByID(id uint32) (*entity.Product, error) {
	var product entity.Product
	err := db.DB().Where("id = ?", id).First(&product).Error
	if product.Id == 0 {
		return nil, nil
	}
	return &product, err
}

func GetProductByName(name string) (*entity.Product, error) {
	var product entity.Product
	err := db.DB().Where("name = ?", name).First(&product).Error
	if product.Id == 0 {
		return nil, nil
	}
	return &product, err
}
