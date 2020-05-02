package recipe

import "go.mongodb.org/mongo-driver/bson/primitive"

// Description is a short view of the recipe
type Description struct {
	ID    *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string              `json:"name"`
	Image bool                `json:"image"`
}
