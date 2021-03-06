package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollection wraps a connected mongo collection
type MongoCollection struct {
	MongoCollection *mongo.Collection
}

// ICollection connects to db
type ICollection interface {
	InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error)
	InsertOneAndFind(ctx context.Context, document interface{}, output interface{}) (interface{}, error)
	DeleteByID(ctx context.Context, objID *primitive.ObjectID) error
	UpdateByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) error
	Find(ctx context.Context, filter interface{}, findOptions *options.FindOptions, obj interface{}) (<-chan interface{}, error)
	FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error)
	FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error)
}

func mongoCollectionIsAnICollection() {
	func(coll ICollection) {}(&MongoCollection{})
}

// CreateCollection gets a wrapped reference to a mongo collection
func CreateCollection(ctx context.Context, database string, collection string) (bool, ICollection) {

	if override, ok := tryGetOverrideFor(database, collection); ok {
		return true, override
	}

	if connected := Connect(ctx); !connected {
		return false, nil
	}

	return true, &MongoCollection{mongoClient.Database(database).Collection(collection)}
}

// UserView returns an ICollection for the mongodb collection users
func UserView(ctx context.Context) (bool, ICollection) {
	return CreateCollection(ctx, "recipeapi", "userviews")
}

// Recipe returns an ICollection for the mongodb collection recipe
func Recipe(ctx context.Context) (bool, ICollection) {
	return CreateCollection(ctx, "recipeapi", "recipes")
}

// Ingredient returns an ICollection for the mongodb collection ingredient
func Ingredient(ctx context.Context) (bool, ICollection) {
	return CreateCollection(ctx, "recipeapi", "ingredients")
}

// Plan returns an ICollection for the mongodb collection plans
func Plan(ctx context.Context) (bool, ICollection) {
	return CreateCollection(ctx, "recipeapi", "plans")
}
