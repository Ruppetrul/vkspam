package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
	"vkspam/models"
	"vkspam/repositories"
)

type UserService interface {
	TryLogin(email string, password string) (string, error)
	CheckEmailExist(email string) (bool, error)
	Register(email string, password string) (*models.User, *string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) TryLogin(email string, password string) (string, error) {
	user, err := s.repo.TryLogin(email)
	if err != nil {
		return "", err
	}
	if user.Email != email {
		return "", errors.New("user with this email not found")
	}
	fmt.Println(user)
	if !checkPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := generateJwtToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

type TokenData struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	t     *jwt.Token
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func generateJwtToken(user *models.User) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	if len(jwtKey) < 1 {
		return "", errors.New("JWT not configured")
	}

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"id":    user.Id,
			"email": user.Email,
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
			"iat":   time.Now().Unix(),
			"iss":   os.Getenv("APP_NAME"),
			"sub":   "user-auth",
			"aud":   []string{"your-audience"},
		})
	return t.SignedString([]byte(jwtKey))
}

func (s *userService) CheckEmailExist(email string) (bool, error) {
	user, err := s.repo.FindUserByEmail(email)

	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	if user.Email == email {
		log.Println(user)
		return true, nil
	}

	return false, nil
}

func (s *userService) Register(email string, password string) (*models.User, *string, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	err = s.repo.Save(&models.User{
		Email:    email,
		Password: passwordHash,
	})
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	token, err := generateJwtToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}
