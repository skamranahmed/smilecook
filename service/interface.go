package service

import (
	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService defines the methods that can be performed on the user object in the service layer
type UserService interface {
	Create(user *models.User) error
	FindOne(username string) (*models.User, error)
	DoesUsernameAlreadyExist(username string) (bool, error)
	HashPassword(plainTextPassword string) (string, error)
	VerifyPassword(plainTextPassword, hashedPassword string) error
}

// RecipeService defines the methods that can be performed on the recipe object in the service layer
type RecipeService interface {
	Create(recipe *models.Recipe) error
	FetchAll() ([]*models.Recipe, error)
	Update(documentObjectID primitive.ObjectID, recipe *models.Recipe) (bool, error)
}
