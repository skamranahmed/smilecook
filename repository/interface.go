package repository

import (
	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository : defines the methods that can be performed on the user object in the repository layer
type UserRepository interface {
	Create(user *models.User) error
	FindOne(username string) (*models.User, error)
	DoesUsernameAlreadyExist(username string) (bool, error)
}

// RecipeRepository : defines the methods that can be performed on the recipe object in the repository layer
type RecipeRepository interface {
	Create(recipe *models.Recipe) error
	FindOne(documentObjectID primitive.ObjectID) (*models.Recipe, error)
	FetchAll() ([]*models.Recipe, error)
	Update(documentObjectID primitive.ObjectID, recipe *models.Recipe) (bool, error)
	Delete(documentObjectID primitive.ObjectID) (bool, error)
}
