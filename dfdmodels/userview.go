package dfdmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserView represents a users set of lists
type UserView struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Nickname string             `json:"nickname"`
	Tags     []string           `json:"tags"`
}
