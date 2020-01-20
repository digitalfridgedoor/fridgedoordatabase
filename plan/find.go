package plan

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByMonthAndYear finds users plan, or creates one
func FindByMonthAndYear(ctx context.Context, userID primitive.ObjectID, month int, year int) (*Plan, error) {
	plan, _, err := getOrCreateOne(ctx, userID, month, year)
	return plan, err
}

// FindOne finds a Plan by id
func FindOne(ctx context.Context, planID primitive.ObjectID) (*Plan, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	singleResult, err := collection.FindByID(ctx, planID.Hex())
	if err != nil {
		return nil, err
	}

	plan, err := fridgedoordatabase.ParseSingleResult(singleResult, &Plan{})
	if err != nil {
		return nil, err
	}

	return plan.(*Plan), nil
}

func getOrCreateOne(ctx context.Context, userID primitive.ObjectID, month int, year int) (*Plan, bool, error) {
	plan, err := findByMonthAndYear(ctx, userID, month, year)
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

func findByMonthAndYear(ctx context.Context, userID primitive.ObjectID, month int, year int) ([]*Plan, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	planBson := bson.M{"month": month, "year": year, "userid": userID}

	cur, err := mongoCollection.Find(ctx, planBson, findOptions)
	if err != nil {
		return nil, err
	}

	planCh := fridgedoordatabase.Parse(ctx, cur, &Plan{})

	results := make([]*Plan, 0)

	for i := range planCh {
		results = append(results, i.(*Plan))
	}

	return results, nil
}
