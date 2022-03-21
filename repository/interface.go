package repository

import "github.com/skamranahmed/smilecook/models"

// UserRepository defines the methods that can be performed on the user object in the repository layer
type UserRepository interface {
	Create(user *models.User) error
	DoesUsernameAlreadyExist(username string) (bool, error)
}
