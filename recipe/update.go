package recipe

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rename changes the name of a recipe
func Rename(ctx context.Context, user primitive.ObjectID, recipeID string, name string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	recipe.Name = name

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// SetImageFlag indicates whether the recipe has an image to look for
func SetImageFlag(ctx context.Context, user primitive.ObjectID, recipeID string, hasImage bool) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	recipe.Image = hasImage

	return collection.UpdateByID(ctx, recipeID, recipe)
}
