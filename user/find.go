package user

import (
	"context"
	"digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a User by ID
func (conn *Connection) FindOne(ctx context.Context, id string) (*User, error) {

	collection := conn.collection()

	findOneOptions := options.FindOne()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	singleResult := collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, findOneOptions)

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
