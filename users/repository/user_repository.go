package repository

import (
	"github.com/alexispell/my-tg/users/model"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByNickname(nickname string) (*model.User, error)
	GetUserByDeviceFingerPrint(deviceFingerPrint string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByNickname(nickname string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("nickname = ?", nickname).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByDeviceFingerPrint(deviceFingerPrint string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("device_fingerprint = ?", deviceFingerPrint).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
