package dfdmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plan represents a meal plan for a month
type Plan struct {
	ID     *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Month  int                 `json:"month"`
	Year   int                 `json:"year"`
	UserID primitive.ObjectID  `json:"userID"`
	Days   []Day               `json:"Days"`
}

// Day represents a meal plan for a day
type Day struct {
	Meal []Meal `json:"meal"`
}

// Meal represents a one meal plan
type Meal struct {
	RecipeID primitive.ObjectID `json:"recipeId"`
	Name     string             `json:"name"`
}
