package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

	return conn.updateByID(ctx, user)
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

	return conn.updateByID(ctx, user)
}

func (conn *Connection) updateByID(ctx context.Context, user *User) error {

	o := options.FindOneAndReplace()

	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}

	singleResult := conn.collection().FindOneAndReplace(ctx, filter, user, o)

	return singleResult.Err()
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
