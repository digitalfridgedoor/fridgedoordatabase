package userview

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetNickname updates the users nickname
func SetNickname(ctx context.Context, view *View, nickname string) error {

	if nickname == "" || view.Nickname == nickname {
		return nil
	}

	view.Nickname = nickname

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	return collection.UpdateByID(ctx, view.ID.Hex(), view)
}

// AddRecipe adds recipe to users list
func AddRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	if _, ok := view.Collections[collectionName]; !ok {
		view.Collections[collectionName] = &RecipeCollection{Name: collectionName}
	}

	view.Collections[collectionName].addRecipe(recipeID)

	return collection.UpdateByID(ctx, viewID, view)
}

// RemoveRecipe removes recipe from users list
func RemoveRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {
	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	if viewCollection, ok := view.Collections[collectionName]; !ok {
		viewCollection.Recipes = fridgedoordatabase.Filter(viewCollection.Recipes, filterFn)
	}

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	return collection.UpdateByID(ctx, viewID, view)
}

// AddTag adds a tag to users list if it isn't already there
func AddTag(ctx context.Context, viewID string, tag string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	hasTag := false
	for _, t := range view.Tags {
		if t == tag {
			hasTag = true
		}
	}

	if !hasTag {
		view.Tags = append(view.Tags, tag)
	}

	return collection.UpdateByID(ctx, viewID, view)
}

// RemoveTag removes a tag from a users list
func RemoveTag(ctx context.Context, viewID string, tag string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	view.Tags = filterTags(view.Tags, tag)

	return collection.UpdateByID(ctx, viewID, view)
}

func (r *RecipeCollection) addRecipe(recipeID primitive.ObjectID) {
	r.Recipes = append(r.Recipes, recipeID)
}

func filterTags(tags []string, tagToRemove string) []string {
	filtered := []string{}

	for _, tag := range tags {
		if tag != tagToRemove {
			filtered = append(filtered, tag)
		}
	}

	return filtered
}
