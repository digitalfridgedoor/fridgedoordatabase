package plan

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plan represents a meal plan for a month
type Plan struct {
	ID     primitive.ObjectID `son:"id" bson:"_id,omitempty"`
	Month  int                `json:"month"`
	Year   int                `json:"year"`
	UserID primitive.ObjectID `json:"userID"`
	Plan   []Day              `json:"plan"`
}

// Day represents a meal plan for a day
type Day struct {
	Date int    `json:"date"`
	Meal []Meal `json:"meal"`
}

// Meal represents a one meal plan
type Meal struct {
	Date int `json:"date"`
}
