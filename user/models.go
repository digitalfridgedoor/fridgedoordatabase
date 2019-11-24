package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a recipe
type User struct {
	ID       primitive.ObjectID   `son:"id" bson:"_id,omitempty"`
	Username string               `json:"username"`
	Recipes  []primitive.ObjectID `json:"recipes"`
}
