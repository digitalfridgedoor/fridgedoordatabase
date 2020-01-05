package userview

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecipeCollection is a users collection of recipes
type RecipeCollection struct {
	Name    string               `json:"name"`
	Recipes []primitive.ObjectID `json:"recipes"`
}

// View represents a users set of lists
type View struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	Username    string                       `json:"username"`
	Nickname    string                       `json:"nickname"`
	Collections map[string]*RecipeCollection `json:"collections"`
	Tags        []string                     `json:"tags"`
}
