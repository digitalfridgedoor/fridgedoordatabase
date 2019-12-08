package recipe

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe represents a recipe
type Recipe struct {
	ID      primitive.ObjectID `son:"id" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	AddedOn time.Time          `json:"addedOn"`
	AddedBy primitive.ObjectID `json:"addedBy"`
	Method  []MethodStep       `json:"method"`
	Recipes []SubRecipe        `json:"recipes"`
}

// MethodStep is an instruction with a collection of ingredients
type MethodStep struct {
	Action      string       `json:"action"`
	Description string       `json:"description"`
	Time        string       `json:"time"`
	Ingredients []Ingredient `json:"ingredients"`
}

// Ingredient is the ingredient linked to each recipe
type Ingredient struct {
	Name         string `json:"name"`
	Amount       string `json:"amount"`
	Preperation  string `json:"preperation"`
	IngredientID string `json:"ingredientId"`
}

// SubRecipe is a pointer to a recipe that makes up the main recipe
type SubRecipe struct {
	Name     string `json:"name"`
	RecipeID string `json:"recipeId"`
}

// Description is a short representation of a recipe
type Description struct {
	ID   primitive.ObjectID `son:"id" bson:"_id,omitempty"`
	Name string             `json:"name"`
}
