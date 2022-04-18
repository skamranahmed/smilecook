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
	if !ur.isCollectionNameCorrect() {
		// TODO: setup custom errors
		return errors.New("incorrect collection name")
	}

	_, err := ur.collection.InsertOne(ur.ctx, u)
	return err
}

// FindOne : finds a user record with the provided username
func (ur *userRepo) FindOne(username string) (*models.User, error) {
	if !ur.isCollectionNameCorrect() {
		return nil, errors.New("incorrect collection name")
	}

	cur := ur.collection.FindOne(ur.ctx, bson.M{"username": username})
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	var user models.User
	err := cur.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DoesUsernameAlreadyExist: checks whether a user with the provided username exists or not
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

// isCollectionNameCorrect : verifies the collection name for the user queries
func (ur *userRepo) isCollectionNameCorrect() bool {
	return ur.collection.Name() == userCollectionName
}
