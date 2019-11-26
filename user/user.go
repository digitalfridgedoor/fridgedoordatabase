package user

import (
	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/mongo"
)

// Collection is a user-wrapped collection
type Collection struct {
	collection *fridgedoordatabase.Collection
}

// New creates an instance of user.Collection
func New(db fridgedoordatabase.Connection) *Collection {
	return &Collection{db.Collection("recipeapi", "users")}
}

func (coll *Collection) mongoCollection() *mongo.Collection {
	return coll.collection.MongoCollection
}
