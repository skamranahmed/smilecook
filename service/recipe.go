package service

import (
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewRecipeService(recipeRepo repository.RecipeRepository) RecipeService {
	return &recipeService{
		recipeRepo: recipeRepo,
	}
}

type recipeService struct {
	recipeRepo repository.RecipeRepository
}

func (rs *recipeService) Create(r *models.Recipe) error {
	return rs.recipeRepo.Create(r)
}

func (rs *recipeService) FetchAll() ([]*models.Recipe, error) {
	return rs.recipeRepo.FetchAll()
}

func (rs *recipeService) Update(documentObjectID primitive.ObjectID, recipe *models.Recipe) (bool, error) {
	return rs.recipeRepo.Update(documentObjectID, recipe)
}
