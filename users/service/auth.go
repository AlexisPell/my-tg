package service

import (
	"errors"
	"time"

	"github.com/alexispell/my-tg/users/model"
	"github.com/alexispell/my-tg/users/repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: *userRepo}
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateToken(email string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) Register(user *model.User) (string, string, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	user.Password = string(hashedPassword)

	// create user in database
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return "", "", err
	}

	// Генерируем токены
	accessToken, err := generateToken(user.Email, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(user.Email, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// generate tokens
	accessToken, err := generateToken(user.Email, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(user.Email, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	// validate refreshToken
	claims, err := validateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// generate new tokens
	accessToken, err := generateToken(claims.Email, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := generateToken(claims.Email, time.Hour*24)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
