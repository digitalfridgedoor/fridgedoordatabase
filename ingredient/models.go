package ingredient

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ingredient represents a node in the ingredient tree
type Ingredient struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	AddedOn  time.Time          `json:"addedOn"`
	ParentID primitive.ObjectID `json:"parentId"`
}
