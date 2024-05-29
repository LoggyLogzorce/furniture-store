package api

import (
	"furniture_store/db"
	"furniture_store/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()

	tests := []struct {
		name     string
		userData map[string]string
		expected bool
	}{
		{
			name: "User does not exist, creation successful",
			userData: map[string]string{
				"username": "testuser",
				"login":    "testlogin1",
				"password": "testpassword",
			},
			expected: true,
		},
		{
			name: "User does not exist, creation successful",
			userData: map[string]string{
				"username": "testuser",
				"login":    "testlogin2",
				"password": "testpassword",
			},
			expected: true,
		},
		{
			name: "User already exists, creation fails",
			userData: map[string]string{
				"username": "testuser",
				"login":    "existinglogin",
				"password": "testpassword",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.DB().Create(&entity.User{Login: "existinglogin"})
			result := CreateUser(tt.userData)
			assert.Equal(t, tt.expected, result)

			// Cleanup database for next test
			db.DB().Exec("TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE")
		})
	}
}

func TestUserRead(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()
	// Подготовка тестовых данных
	testUserData := map[string]string{
		"login":    "testlogin",
		"password": "testpassword",
	}

	expectedRole := "user"

	// Создание тестового пользователя в тестовой базе данных
	testUser := entity.User{
		Login:    "testlogin",
		Password: "testpassword",
		Role:     "user",
	}
	err := db.DB().Create(&testUser).Error
	assert.NoError(t, err)

	// Вызов тестируемой функции
	cookie, role := UserRead(testUserData)

	// Проверка результата
	assert.NotEmpty(t, cookie.Value)
	assert.Equal(t, expectedRole, role)

	db.DB().Exec("TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE")
}

func TestDeleteUser(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()
	// Подготовка тестовых данных
	testRowData := map[string]string{
		"uid": "1",
	}

	// Создание тестового пользователя
	testUser := entity.User{
		Uid: 1,
		// Здесь можно указать остальные поля пользователя, если они важны для теста
	}

	err := db.DB().Create(&testUser).Error
	assert.NoError(t, err)

	// Вызов тестируемой функции
	DeleteUser(testRowData)

	// Проверка, что пользователь удален
	deletedUser, err := GetUserByID(testUser.Uid)
	if err != nil {
		t.Errorf("Пользователь найден: %v", err)
	}
	assert.Nil(t, deletedUser, "Пользователь не должен быть найден после удаления")
	db.DB().Exec("TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE")
}

func TestUpdateUser(t *testing.T) {
	db.ConnectToTestDB()
	db.Migrate()
	// Подготовка тестовых данных
	updatedData := map[string]string{
		"uid":    "1",
		"Имя":    "Новое имя",
		"Логин":  "новый_логин",
		"Пароль": "новый_пароль",
		"Роль":   "новая_роль",
	}

	// Создание тестового пользователя в базе данных
	originalUser := entity.User{
		Uid:      1,
		Name:     "Имя",
		Login:    "логин",
		Password: "пароль",
		Role:     "роль",
	}
	err := db.DB().Create(&originalUser).Error
	assert.NoError(t, err)

	// Вызов тестируемой функции
	UpdateUser(updatedData)

	// Получение обновленного пользователя из базы данных
	updatedUser, err := GetUserByID(originalUser.Uid)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser, "Пользователь должен быть найден в базе данных")

	// Проверка, что данные пользователя обновлены
	assert.Equal(t, updatedData["Имя"], updatedUser.Name)
	assert.Equal(t, updatedData["Логин"], updatedUser.Login)
	assert.Equal(t, updatedData["Пароль"], updatedUser.Password)
	assert.Equal(t, updatedData["Роль"], updatedUser.Role)
	db.DB().Exec("TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE")
}

func GetUserByID(id uint32) (*entity.User, error) {
	var user entity.User
	err := db.DB().Where("uid = ?", id).First(&user).Error
	if user.Uid == 0 {
		return nil, nil
	}
	return &user, err
}
