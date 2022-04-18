package service

import (
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewRecipeService : returns a recipeService struct that implements the RecipeService interface
func NewRecipeService(recipeRepo repository.RecipeRepository) RecipeService {
	return &recipeService{
		recipeRepo: recipeRepo,
	}
}

type recipeService struct {
	recipeRepo repository.RecipeRepository
}

// Create : creates a new recipe record
func (rs *recipeService) Create(r *models.Recipe) error {
	return rs.recipeRepo.Create(r)
}

// FindOne : finds a user record with the provided ID
func (rs *recipeService) FindOne(documentObjectID primitive.ObjectID) (*models.Recipe, error) {
	return rs.recipeRepo.FindOne(documentObjectID)
}

// FetchAll : fetched all the recipe records
func (rs *recipeService) FetchAll() ([]*models.Recipe, error) {
	return rs.recipeRepo.FetchAll()
}

// Update : updates a recipe record with the provided ID
func (rs *recipeService) Update(documentObjectID primitive.ObjectID, recipe *models.Recipe) (bool, error) {
	return rs.recipeRepo.Update(documentObjectID, recipe)
}

// Delete : deletes a recipe record with the provided ID
func (rs *recipeService) Delete(documentObjectID primitive.ObjectID) (bool, error) {
	return rs.recipeRepo.Delete(documentObjectID)
}
