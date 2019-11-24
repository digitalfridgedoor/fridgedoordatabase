package user

import (
	"digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/mongo"
)

// Connection can find and parse Users from mongodb
type Connection struct {
	db fridgedoordatabase.Connection
}

// New creates an instance of recipe.Connection
func New(db fridgedoordatabase.Connection) *Connection {
	return &Connection{db}
}

func (conn *Connection) collection() *mongo.Collection {
	return conn.db.Collection("recipeapi", "users")
}
