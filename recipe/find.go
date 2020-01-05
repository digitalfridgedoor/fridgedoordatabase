package recipe

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a recipe by ID
func FindOne(ctx context.Context, id string) (*Recipe, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	singleResult, err := collection.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})
	if err != nil {
		return nil, err
	}

	return ing.(*Recipe), err
}

// FindByIds finds recipe by ID
func FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]*Description, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	_in := bson.M{"$in": ids}
	idin := bson.M{"_id": _in}

	// todo: projection to only select the fields in Description?
	cur, err := mongoCollection.Find(context.Background(), idin, findOptions)
	if err != nil {
		return make([]*Description, 0), err
	}

	return parseRecipe(ctx, cur)
}

// FindByName finds recipes starting with the given letter
func FindByName(ctx context.Context, startsWith string, userID primitive.ObjectID) ([]*Recipe, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}
	addedByBson := bson.M{"addedby": userID}
	andBson := bson.M{"$and": []bson.M{startsWithBson, addedByBson}}

	cur, err := mongoCollection.Find(ctx, andBson, findOptions)
	if err != nil {
		return make([]*Recipe, 0), err
	}

	recipeCh := fridgedoordatabase.Parse(ctx, cur, &Recipe{})

	results := make([]*Recipe, 0)

	for i := range recipeCh {
		results = append(results, i.(*Recipe))
	}

	return results, nil
}

// UserRecipes lists all the users recipes
func UserRecipes(ctx context.Context, userID primitive.ObjectID) ([]*Description, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	cur, err := mongoCollection.Find(findctx, bson.D{primitive.E{Key: "addedBy", Value: userID}}, findOptions)
	if err != nil {
		return make([]*Description, 0), err
	}

	return parseRecipe(ctx, cur)
}

// FindByTags finds recipes with the given tags
func FindByTags(ctx context.Context, userID primitive.ObjectID, tags []string, notTags []string) ([]*Recipe, error) {

	// https://stackoverflow.com/questions/6940503/mongodb-get-documents-by-tags

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	// { $and: [ {tags: { $all: ["tag"] } }, { tags: { $nin: ["anothertag"] } } ] }

	addedByBson := bson.M{"addedby": userID}
	andBson := []bson.M{addedByBson}

	if tags != nil && len(tags) > 0 {
		allBson := bson.M{"$all": tags}
		tagsBson := bson.M{"metadata.tags": allBson}
		andBson = append(andBson, tagsBson)
	}

	if notTags != nil && len(notTags) > 0 {
		ninBson := bson.M{"$nin": notTags}
		ninTagsBson := bson.M{"metadata.tags": ninBson}
		andBson = append(andBson, ninTagsBson)
	}

	cur, err := mongoCollection.Find(ctx, bson.M{"$and": andBson}, findOptions)
	if err != nil {
		return make([]*Recipe, 0), err
	}

	recipeCh := fridgedoordatabase.Parse(ctx, cur, &Recipe{})

	results := make([]*Recipe, 0)

	for i := range recipeCh {
		results = append(results, i.(*Recipe))
	}

	return results, nil
}

func parseRecipe(ctx context.Context, cur *mongo.Cursor) ([]*Description, error) {
	ingCh := fridgedoordatabase.Parse(ctx, cur, &Description{})

	results := make([]*Description, 0)

	for i := range ingCh {
		results = append(results, i.(*Description))
	}

	return results, nil
}
