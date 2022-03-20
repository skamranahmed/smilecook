package repository

import (
	"context"
	"errors"

	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollectionName string = "users"

// NewUserRepository : returns a userRepo struct that implements the UserRepository interface
func NewUserRepository(ctx context.Context, userCollection *mongo.Collection) UserRepository {
	return &userRepo{
		ctx:        ctx,
		collection: userCollection,
	}
}

type userRepo struct {
	ctx        context.Context
	collection *mongo.Collection
}

// Create: inserts a new user record in the `users` collection
func (ur *userRepo) Create(u *models.User) error {
	collectionName := ur.collection.Name()
	if collectionName != userCollectionName {
		// TODO: setup custom errors
		return errors.New("incorrect collection name")
	}
	_, err := ur.collection.InsertOne(ur.ctx, u)
	return err
}
