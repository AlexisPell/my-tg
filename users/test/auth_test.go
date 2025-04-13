package test

import (
	"testing"

	"github.com/alexispell/my-tg/users/model"
	"github.com/alexispell/my-tg/users/service"

	"github.com/alexispell/my-tg/users/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(&repo)

	user := &model.User{
		Nickname:          "testuser1",
		Email:             "test1@example.com",
		Password:          "password123",
		DeviceFingerPrint: "unique_device_fingerprint1",
	}
	err := service.Register(user)

	assert.Nil(t, err)
	assert.NotEmpty(t, user.Password)
}

func TestLogin(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(&repo)

	user := &model.User{
		Nickname:          "testuser2",
		Email:             "test2@example.com",
		Password:          "password123",
		DeviceFingerPrint: "unique_device_fingerprint2",
	}
	service.Register(user)

	loggedInUser, err := service.Login("test2@example.com", "password123")

	assert.Nil(t, err)
	assert.Equal(t, user.Email, loggedInUser.Email)
}

func TestLoginWithWrongPassword(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(&repo)

	user := &model.User{
		Nickname:          "testuser3",
		Email:             "test3@example.com",
		Password:          "correctpassword",
		DeviceFingerPrint: "unique_device_fingerprint3",
	}
	service.Register(user)

	loggedInUser, err := service.Login("test3@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, loggedInUser)
}
