package api

import (
	"furniture_store/db"
	"furniture_store/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateCategory(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	updatedData := map[string]string{
		"id":        "1",
		"Категория": "Кресла",
	}

	testData := entity.Category{
		Id:   1,
		Name: "Стулья",
	}

	err := db.DB().Create(&testData).Error
	assert.NoError(t, err)

	UpdateCategory(updatedData)

	var updatedCategory entity.Category
	err = db.DB().Where("id = ?", testData.Id).First(&updatedCategory).Error
	assert.NoError(t, err)
	assert.Equal(t, updatedData["Категория"], updatedCategory.Name)
	db.DB().Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
}

func TestDeleteCategory(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	deleteData := map[string]string{
		"id": "1",
	}

	testData := entity.Category{
		Id:   1,
		Name: "Стулья",
	}

	err := db.DB().Create(&testData).Error
	assert.NoError(t, err)

	DeleteCategory(deleteData)

	deletedData, err := GetCategoryByID(testData.Id)
	assert.NoError(t, err)
	assert.Nil(t, deletedData, "Категория должна быть удалена из базы данных")
	db.DB().Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
}

func TestAddCategory(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	addData := map[string]string{
		"category": "Кресла",
	}

	AddCategory(addData)

	var addedCategory entity.Category
	err := db.DB().Where("name = ?", addData["category"]).First(&addedCategory).Error
	assert.NoError(t, err)
	assert.Equal(t, addData["category"], addedCategory.Name)
	db.DB().Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
}

func GetCategoryByID(id uint32) (*entity.Category, error) {
	var category entity.Category
	err := db.DB().Where("id = ?", id).First(&category).Error
	if category.Id == 0 {
		return nil, nil
	}
	return &category, err
}
