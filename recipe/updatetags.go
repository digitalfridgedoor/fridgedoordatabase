package recipe

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddTag adds a new tag to the recipe
func AddTag(ctx context.Context, user primitive.ObjectID, recipeID string, tag string) error {

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

	recipe.Tags = appendTag(recipe.Tags, tag)

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveTag removes a tag from the recipe
func RemoveTag(ctx context.Context, user primitive.ObjectID, recipeID string, tag string) error {

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

	recipe.Tags = removeTag(recipe.Tags, tag)

	return collection.UpdateByID(ctx, recipeID, recipe)
}

func appendTag(tags []string, tag string) []string {
	hasTag := false

	for _, t := range tags {
		if t == tag {
			hasTag = true
		}
	}

	if !hasTag {
		tags = append(tags, tag)
	}

	return tags
}

func removeTag(tags []string, tag string) []string {
	filtered := []string{}

	for _, t := range tags {
		if t != tag {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
