package recipe

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rename changes the name of a recipe
func Rename(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, name string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID, userID)
	if err != nil {
		return nil, err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return nil, errUnauthorised
	}

	recipe.Name = name

	err = coll.c.UpdateByID(ctx, recipeID, recipe)
	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, recipeID, userID)
}

// UpdateMetadata updates any information in metadata
func UpdateMetadata(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, updates map[string]string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID, userID)
	if err != nil {
		return nil, err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return nil, errUnauthorised
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

	err = coll.c.UpdateByID(ctx, recipeID, recipe)
	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, recipeID, userID)
}
