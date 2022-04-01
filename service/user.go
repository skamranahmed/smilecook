package service

import (
	"fmt"

	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/repository"
	"golang.org/x/crypto/bcrypt"
)

// NewUserService : returns a userService struct that implements the UserService interface
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type userService struct {
	userRepo repository.UserRepository
}

func (us *userService) Create(u *models.User) error {
	return us.userRepo.Create(u)
}

func (us *userService) FindOne(username string) (*models.User, error) {
	return us.userRepo.FindOne(username)
}

func (us *userService) DoesUsernameAlreadyExist(username string) (bool, error) {
	return us.userRepo.DoesUsernameAlreadyExist(username)
}

func (us *userService) HashPassword(plainTextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("falied to hash password, error: %s", err)
	}
	return string(hashedPassword), nil
}

func (us *userService) VerifyPassword(plainTextPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
}
