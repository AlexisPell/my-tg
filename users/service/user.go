package service

import (
	"github.com/alexispell/my-tg/users/model"
	"github.com/alexispell/my-tg/users/repository"
)

// UserService предоставляет методы для работы с пользователями

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.UserRepo.CreateUser(user)
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.UserRepo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.UserRepo.GetUserByEmail(email)
}

func (s *UserService) GetUserByNickname(nickname string) (*model.User, error) {
	return s.UserRepo.GetUserByNickname(nickname)
}

func (s *UserService) GetUserByDeviceFingerPrint(deviceFingerPrint string) (*model.User, error) {
	return s.UserRepo.GetUserByDeviceFingerPrint(deviceFingerPrint)
}
