package repository

import (
	"context"
	"errors"

	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (rr *recipeRepo) FetchAll() ([]*models.Recipe, error) {
	cur, err := rr.collection.Find(rr.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(rr.ctx)

	recipes := make([]*models.Recipe, 0)
	for cur.Next(rr.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

func (rr *recipeRepo) isCollectionNameCorrect() bool {
	return rr.collection.Name() == recipeCollectionName
}
