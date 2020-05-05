package userview

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
