package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddRecipe adds recipe to users list
func (conn *Connection) AddRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {
	user, err := conn.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	user.Recipes = append(user.Recipes, recipeID)

	o := options.Update()

	singleResult := conn.collection().FindOneAndUpdate(ctx, user, o)

	return singleResult.Err()
}

// RemoveRecipe removes recipe from users list
func (conn *Connection) RemoveRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {
	user, err := conn.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	user.Recipes = filter(user.Recipes, filterFn)

	o := options.Update()

	_, err = conn.collection().UpdateOne(ctx, user, o)

	return err
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
