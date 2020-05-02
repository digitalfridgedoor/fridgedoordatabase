package recipe

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a recipe by ID
func FindOne(ctx context.Context, id *primitive.ObjectID, userID primitive.ObjectID) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	return coll.findOne(ctx, id, userID)
}

// FindByName finds recipes starting with the given letter
func FindByName(ctx context.Context, startsWith string, userID primitive.ObjectID, limit int64) ([]*Description, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	if limit > 20 {
		limit = 20
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}
	addedByBson := bson.M{"addedby": userID}
	andBson := bson.M{"$and": []bson.M{startsWithBson, addedByBson}}

	ch, err := coll.c.Find(ctx, andBson, findOptions, &dfdmodels.Recipe{})
	if err != nil {
		return []*Description{}, err
	}

	results := readChannel(ch, userID)
	return results, nil
}

// FindByTags finds recipes with the given tags
func FindByTags(ctx context.Context, userID primitive.ObjectID, tags []string, notTags []string, limit int64) ([]*Description, error) {

	// https://stackoverflow.com/questions/6940503/mongodb-get-documents-by-tags

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	if limit > 20 {
		limit = 20
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)

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

	ch, err := coll.c.Find(ctx, bson.M{"$and": andBson}, findOptions, &dfdmodels.Recipe{})
	if err != nil {
		return []*Description{}, err
	}

	results := readChannel(ch, userID)
	return results, nil
}

// FindPublic gets a users public recipes
func FindPublic(ctx context.Context, userID primitive.ObjectID, limit int64) ([]*Description, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	if limit > 20 {
		limit = 20
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	addedByBson := bson.M{"addedby": userID}
	viewableByEveryone := bson.M{"metadata.viewableby.everyone": true}
	andBson := []bson.M{addedByBson, viewableByEveryone}

	ch, err := coll.c.Find(ctx, bson.M{"$and": andBson}, findOptions, &dfdmodels.Recipe{})
	if err != nil {
		return make([]*Description, 0), err
	}

	results := readChannel(ch, userID)
	return results, nil
}

func (coll *collection) findOne(ctx context.Context, id *primitive.ObjectID, userID primitive.ObjectID) (*dfdmodels.Recipe, error) {

	r, err := coll.c.FindByID(ctx, id, &dfdmodels.Recipe{})

	if err != nil {
		return nil, err
	}

	re := r.(*dfdmodels.Recipe)
	if !CanView(re, userID) {
		return nil, errUnauthorised
	}

	return re, err
}

func readChannel(ch <-chan interface{}, userID primitive.ObjectID) []*Description {
	results := make([]*Description, 0)

	for i := range ch {
		r := i.(*dfdmodels.Recipe)

		if CanView(r, userID) {
			results = append(results, &Description{
				ID:    r.ID,
				Name:  r.Name,
				Image: r.Metadata.Image,
			})
		}
	}

	return results
}
