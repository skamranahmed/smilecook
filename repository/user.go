package repository

import (
	"context"
	"errors"

	"github.com/skamranahmed/smilecook/models"
	"go.mongodb.org/mongo-driver/bson"
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
	// collectionName := ur.collection.Name()
	// if collectionName != userCollectionName {
	// 	// TODO: setup custom errors
	// 	return errors.New("incorrect collection name")
	// }

	if !ur.isCollectionNameCorrect() {
		// TODO: setup custom errors
		return errors.New("incorrect collection name")
	}

	_, err := ur.collection.InsertOne(ur.ctx, u)
	return err
}

// FindOne: checks whether a user with the provided credentials exists or not
func (ur *userRepo) DoesUsernameAlreadyExist(username string) (bool, error) {
	if !ur.isCollectionNameCorrect() {
		return false, errors.New("incorrect collection name")
	}

	cur := ur.collection.FindOne(ur.ctx, bson.M{
		"username": username,
	})

	if cur.Err() != mongo.ErrNoDocuments {
		// this means a user with the provided username already exists
		return true, nil
	}

	return false, nil
}

func (ur *userRepo) isCollectionNameCorrect() bool {
	return ur.collection.Name() == userCollectionName
}
