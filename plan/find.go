package plan

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByMonthAndYear finds users plan, or creates one
func FindByMonthAndYear(ctx context.Context, userID primitive.ObjectID, month int, year int) (*dfdmodels.Plan, error) {
	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	plan, _, err := coll.getOrCreateOne(ctx, userID, month, year)
	return plan, err
}

// FindOne finds a Plan by id
func FindOne(ctx context.Context, planID *primitive.ObjectID) (*dfdmodels.Plan, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	plan, err := coll.c.FindByID(ctx, planID, &dfdmodels.Plan{})
	if err != nil {
		return nil, err
	}

	return plan.(*dfdmodels.Plan), nil
}

func (coll *collection) getOrCreateOne(ctx context.Context, userID primitive.ObjectID, month int, year int) (*dfdmodels.Plan, bool, error) {

	plan, err := coll.findByMonthAndYear(ctx, userID, month, year)
	if err != nil {
		return nil, false, err
	}

	if len(plan) == 0 {
		if ok, p := create(userID, month, year); ok {
			return p, true, nil
		}

		return nil, false, errInvalidInput
	}

	return plan[0], false, nil
}

func (coll *collection) findByMonthAndYear(ctx context.Context, userID primitive.ObjectID, month int, year int) ([]*dfdmodels.Plan, error) {

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	planBson := bson.M{"month": month, "year": year, "userid": userID}

	ch, err := coll.c.Find(ctx, planBson, findOptions, &dfdmodels.Plan{})

	if err != nil {
		return nil, err
	}
	results := make([]*dfdmodels.Plan, 0)

	for i := range ch {
		results = append(results, i.(*dfdmodels.Plan))
	}

	return results, nil
}
