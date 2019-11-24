package recipe

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

// Connection can find and parse Recipes from mongodb
type Connection struct {
	db fridgedoordatabase.Connection
}

// New creates an instance of recipe.Connection
func New(db fridgedoordatabase.Connection) *Connection {
	return &Connection{db}
}

func (conn *Connection) collection() *mongo.Collection {
	return conn.db.Collection("recipeapi", "recipes")
}
