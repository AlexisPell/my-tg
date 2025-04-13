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
	accessToken, refreshToken, err := service.Register(user)

	assert.Nil(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestLogin_Success(t *testing.T) {
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

	accessToken, refreshToken, err := service.Login("test2@example.com", "password123")

	assert.Nil(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestLogin_Failure_WithWrongPassword(t *testing.T) {
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

	accessToken, refreshToken, err := service.Login("test3@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
}

func TestRefreshTokens_Success(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(&repo)

	user := &model.User{
		Nickname:          "testuser4",
		Email:             "test4@example.com",
		Password:          "password123",
		DeviceFingerPrint: "unique_device_fingerprint4",
	}
	_, refreshToken, _ := service.Register(user)

	newAccessToken, newRefreshToken, err := service.RefreshTokens(refreshToken)

	assert.Nil(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
}

func TestRefreshTokens_Failure(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(&repo)

	invalidRefreshToken := "invalidToken"

	newAccessToken, newRefreshToken, err := service.RefreshTokens(invalidRefreshToken)

	assert.NotNil(t, err)
	assert.Empty(t, newAccessToken)
	assert.Empty(t, newRefreshToken)
}
