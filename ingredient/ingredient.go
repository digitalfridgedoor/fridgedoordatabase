package ingredient

import (
	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/mongo"
)

// Collection is a ingredient-wrapped collection
type Collection struct {
	collection *fridgedoordatabase.Collection
}

// New creates an instance of ingredient.Collection
func New(db fridgedoordatabase.Connection) *Collection {
	return &Collection{db.Collection("recipeapi", "ingredients")}
}

func (coll *Collection) mongoCollection() *mongo.Collection {
	return coll.collection.MongoCollection
}
