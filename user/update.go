package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddRecipe adds recipe to users list
func (coll *Collection) AddRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {

	user, err := coll.MongoCollection.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	user.Recipes = append(user.Recipes, recipeID)

	return coll.UpdateByID(ctx, user)
}

// RemoveRecipe removes recipe from users list
func (coll *Collection) RemoveRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {
	user, err := coll.MongoCollection.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	user.Recipes = filter(user.Recipes, filterFn)

	return coll.UpdateByID(ctx, user)
}

func filter(ids []primitive.ObjectID, filterFn func(id *primitive.ObjectID) bool) []primitive.ObjectID {
	filtered := []primitive.ObjectID{}

	for id := range iterateObjectIDs(ids) {
		if filterFn(id) {
			filtered = append(filtered, *id)
		}
	}

	return filtered
}

func iterateObjectIDs(ids []primitive.ObjectID) <-chan *primitive.ObjectID {
	ch := make(chan *primitive.ObjectID)

	go func() {
		defer close(ch)
		for _, id := range ids {
			ch <- &id
		}
	}()

	return ch
}
