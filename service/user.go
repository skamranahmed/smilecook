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

// Create : creates a new user record
func (us *userService) Create(u *models.User) error {
	return us.userRepo.Create(u)
}

// FindOne : finds a user record with the provided username
func (us *userService) FindOne(username string) (*models.User, error) {
	return us.userRepo.FindOne(username)
}

// DoesUsernameAlreadyExist: checks whether a user with the provided username exists or not
func (us *userService) DoesUsernameAlreadyExist(username string) (bool, error) {
	return us.userRepo.DoesUsernameAlreadyExist(username)
}

// HashPassword : hashes a plainTextPassword
func (us *userService) HashPassword(plainTextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("falied to hash password, error: %s", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword : verifies whether the hash of the provided plainTextPassword matches with the existing hashPassword or not
func (us *userService) VerifyPassword(plainTextPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
}
