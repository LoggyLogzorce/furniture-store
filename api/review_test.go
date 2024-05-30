package api

import (
	"furniture_store/db"
	"furniture_store/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteReview(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()
	// Arrange
	testData := map[string]string{
		"id": "1",
	}

	testReview := entity.Review{Id: 1}

	err := db.DB().Create(&testReview).Error
	assert.NoError(t, err)

	DeleteReview(testData)

	deletedReview, err := GetReviewByID(testReview.Id)
	assert.NoError(t, err)
	assert.Nil(t, deletedReview, "Отзыв должен быть удален из базы данных")
	db.DB().Exec("TRUNCATE TABLE review RESTART IDENTITY CASCADE")
}

func GetReviewByID(id uint32) (*entity.Review, error) {
	var review entity.Review
	err := db.DB().Where("id = ?", id).First(&review).Error
	if review.Id == 0 {
		return nil, nil
	}
	return &review, err
}
