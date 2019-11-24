package recipe

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe represents a recipe
type Recipe struct {
	ID          primitive.ObjectID `son:"id" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	AddedOn     time.Time          `json:"addedOn"`
	Ingredients []Ingredient       `json:"ingredients"`
	Recipes     []SubRecipe        `json:"recipes"`
}

// Ingredient is the ingredient linked to each recipe
type Ingredient struct {
	Name         string `json:"name"`
	Amount       string `json:"amount"`
	IngredientID string `json:"ingredientId"`
}

// SubRecipe is a pointer to a recipe that makes up the main recipe
type SubRecipe struct {
	Name     string `json:"name"`
	RecipeID string `json:"recipeId"`
}
