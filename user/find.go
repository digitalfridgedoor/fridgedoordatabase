package user

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a User by ID
func (coll *Collection) FindOne(ctx context.Context, id string) (*User, error) {

	singleResult := coll.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &User{})
	if err != nil {
		return nil, err
	}

	return ing.(*User), err
}

// GetByUsername tries to get User by username
func (conn *Connection) GetByUsername(ctx context.Context, username string) (*User, error) {

	collection := conn.collection()

	// Pass these options to the FindOne method
	findOneOptions := options.FindOne()

	singleResult := collection.FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}, findOneOptions)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &User{})
	if err != nil {
		return nil, err
	}

	return ing.(*User), err
}
