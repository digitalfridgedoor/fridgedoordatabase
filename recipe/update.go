package recipe

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rename changes the name of a recipe
func Rename(ctx context.Context, user primitive.ObjectID, recipeID *primitive.ObjectID, name string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !CanEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	recipe.Name = name

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

// UpdateMetadata updates any information in metadata
func UpdateMetadata(ctx context.Context, user primitive.ObjectID, recipeID *primitive.ObjectID, updates map[string]string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !CanEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	if update, ok := updates["image"]; ok {
		recipe.Metadata.Image = update == "true"
	}
	if update, ok := updates["tag_add"]; ok {
		recipe.Metadata.Tags = appendString(recipe.Metadata.Tags, update)
	}
	if update, ok := updates["tag_remove"]; ok {
		recipe.Metadata.Tags = removeString(recipe.Metadata.Tags, update)
	}
	if update, ok := updates["link_add"]; ok {
		recipe.Metadata.Links = appendString(recipe.Metadata.Links, update)
	}
	if update, ok := updates["link_remove"]; ok {
		recipe.Metadata.Links = removeString(recipe.Metadata.Links, update)
	}
	if update, ok := updates["viewableby_everyone"]; ok {
		recipe.Metadata.ViewableBy.Everyone = update == "true"
	}
	if update, ok := updates["viewableby_adduser"]; ok {
		objectID, err := primitive.ObjectIDFromHex(update)
		if err == nil {
			recipe.Metadata.ViewableBy.Users = appendID(recipe.Metadata.ViewableBy.Users, objectID)
		}
	}
	if update, ok := updates["viewableby_removeuser"]; ok {
		objectID, err := primitive.ObjectIDFromHex(update)
		if err == nil {
			recipe.Metadata.ViewableBy.Users = removeID(recipe.Metadata.ViewableBy.Users, objectID)
		}
	}

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}
