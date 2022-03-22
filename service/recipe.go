package service

import (
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/repository"
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
