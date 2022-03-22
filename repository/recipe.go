package repository

import (
	"context"
	"errors"

	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const recipeCollectionName string = "recipes"

func NewRecipeRepository(ctx context.Context, recipeCollection *mongo.Collection) RecipeRepository {
	return &recipeRepo{
		ctx:        ctx,
		collection: recipeCollection,
	}
}

type recipeRepo struct {
	ctx        context.Context
	collection *mongo.Collection
}

// Create: inserts a new recipe record in the `recipes` collection
func (rr *recipeRepo) Create(r *models.Recipe) error {
	if !rr.isCollectionNameCorrect() {
		// TODO: setup custom errors
		return errors.New("incorrect collection name")
	}

	_, err := rr.collection.InsertOne(rr.ctx, r)
	return err
}

func (rr *recipeRepo) isCollectionNameCorrect() bool {
	return rr.collection.Name() == recipeCollectionName
}
