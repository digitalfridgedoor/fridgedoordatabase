package recipe

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

// Collection is a recipe-wrapped collection
type Collection struct {
	collection *fridgedoordatabase.Collection
}

// New creates an instance of recipe.Collection
func New(db fridgedoordatabase.Connection) *Collection {
	return &Collection{db.Collection("recipeapi", "recipes")}
}

func (coll *Collection) mongoCollection() *mongo.Collection {
	return coll.collection.MongoCollection
}
