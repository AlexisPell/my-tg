package service

import (
	"errors"

	"github.com/alexispell/my-tg/users/model"
	"github.com/alexispell/my-tg/users/repository"

	"golang.org/x/crypto/bcrypt"
)

// AuthService предоставляет методы для аутентификации пользователей

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: *userRepo}
}

func (s *AuthService) Register(user *model.User) error {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Создаем пользователя в базе данных
	return s.UserRepo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
