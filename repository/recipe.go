package repository

import (
	"context"
	"errors"

	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const recipeCollectionName string = "recipes"

// NewRecipeRepository : returns a recipeRepo struct that implements the RecipeRepository interface
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
		return errors.New("incorrect collection name")
	}

	_, err := rr.collection.InsertOne(rr.ctx, r)
	return err
}

// FindOne : finds a recipe record with the provided id
func (rr *recipeRepo) FindOne(documentObjectID primitive.ObjectID) (*models.Recipe, error) {
	if !rr.isCollectionNameCorrect() {
		return nil, errors.New("incorrect collection name")
	}

	cur := rr.collection.FindOne(rr.ctx, bson.M{"_id": documentObjectID})

	var recipe models.Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

// FetchAll : fetches all public recipe records
func (rr *recipeRepo) FetchAll() ([]*models.Recipe, error) {
	if !rr.isCollectionNameCorrect() {
		return nil, errors.New("incorrect collection name")
	}

	cur, err := rr.collection.Find(rr.ctx, bson.M{"isPrivate": false})
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

// Update : updates a recipe record with the provided ID
func (rr *recipeRepo) Update(documentObjectID primitive.ObjectID, recipe *models.Recipe) (bool, error) {
	if !rr.isCollectionNameCorrect() {
		return false, errors.New("incorrect collection name")
	}

	result, err := rr.collection.UpdateOne(rr.ctx,
		bson.M{"_id": documentObjectID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: recipe.Name},
				{Key: "instructions", Value: recipe.Instructions},
				{Key: "ingredients", Value: recipe.Ingredients},
				{Key: "tags", Value: recipe.Tags},
			},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, nil
	}

	return true, nil
}

// Delete : deletes a recipe record with the provided ID
func (rr *recipeRepo) Delete(documentObjectID primitive.ObjectID) (bool, error) {
	if !rr.isCollectionNameCorrect() {
		return false, errors.New("incorrect collection name")
	}

	result, err := rr.collection.DeleteOne(rr.ctx, bson.M{"_id": documentObjectID})
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}

// isCollectionNameCorrect : verifies the collection name for the recipe queries
func (rr *recipeRepo) isCollectionNameCorrect() bool {
	return rr.collection.Name() == recipeCollectionName
}
