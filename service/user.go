package service

import (
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/repository"
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
