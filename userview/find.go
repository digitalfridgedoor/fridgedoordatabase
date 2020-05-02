package userview

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a User by ID
func FindOne(ctx context.Context, id *primitive.ObjectID) (*dfdmodels.UserView, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	return coll.findOne(ctx, id)
}

// GetByUsername tries to get User by username
func GetByUsername(ctx context.Context, username string) (*dfdmodels.UserView, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	return coll.getByUsername(ctx, username)
}

// GetLinkedUserViews returns all user views for now
func GetLinkedUserViews(ctx context.Context) ([]*dfdmodels.UserView, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	ch, err := coll.c.Find(findctx, bson.D{{}}, findOptions, &dfdmodels.UserView{})
	if err != nil {
		return nil, err
	}

	results := make([]*dfdmodels.UserView, 0)
	for i := range ch {
		results = append(results, i.(*dfdmodels.UserView))
	}
	return results, nil
}

func (coll *collection) findOne(ctx context.Context, id *primitive.ObjectID) (*dfdmodels.UserView, error) {

	uv, err := coll.c.FindByID(ctx, id, &dfdmodels.UserView{})
	if err != nil {
		return nil, err
	}

	return uv.(*dfdmodels.UserView), err
}

func (coll *collection) getByUsername(ctx context.Context, username string) (*dfdmodels.UserView, error) {

	uv, err := coll.c.FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}, &dfdmodels.UserView{})

	if err != nil {
		return nil, err
	}

	return uv.(*dfdmodels.UserView), nil
}
