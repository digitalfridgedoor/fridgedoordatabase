package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddRecipeToUser adds recipe to user
func (conn *Connection) AddRecipeToUser(ctx context.Context, userID string, recipeID primitive.ObjectID) error {
	user, err := conn.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	user.Recipes = append(user.Recipes, recipeID)

	o := options.Update()

	_, err = conn.collection().UpdateOne(ctx, user, o)

	return err
}
