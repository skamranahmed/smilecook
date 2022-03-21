package service

import "github.com/skamranahmed/smilecook/models"

// UserService defines the methods that can be performed on the user object in the service layer
type UserService interface {
	Create(user *models.User) error
	DoesUsernameAlreadyExist(username string) (bool, error)
}
