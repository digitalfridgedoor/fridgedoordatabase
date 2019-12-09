package ingredient

import (
	"context"
	"log"

	"github.com/digitalfridgedoor/fridgedoordatabase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByName finds ingredients starting with the given letter
func (coll *Collection) FindByName(ctx context.Context, startsWith string) ([]*Ingredient, error) {

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}

	cur, err := coll.mongoCollection().Find(ctx, startsWithBson, findOptions)
	if err != nil {
		return make([]*Ingredient, 0), err
	}

	ingCh := fridgedoordatabase.Parse(ctx, cur, &Ingredient{})

	results := make([]*Ingredient, 0)

	for i := range ingCh {
		results = append(results, i.(*Ingredient))
	}

	return results, nil
}

// FindOne does not find one
func (coll *Collection) FindOne(ctx context.Context, id string) (*Ingredient, error) {

	singleResult, err := coll.collection.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Ingredient{})

	return ing.(*Ingredient), err
}

// IngredientByParentID returns an array of ingredients with the given parentID
func (coll *Collection) IngredientByParentID(ctx context.Context, parentID primitive.ObjectID) []*Ingredient {

	cur, err := coll.mongoCollection().Find(ctx, bson.M{"parentId": parentID})
	if err != nil {
		log.Fatal(err)
	}

	ingCh := fridgedoordatabase.Parse(ctx, cur, &Ingredient{})

	results := make([]*Ingredient, 0)

	for i := range ingCh {
		results = append(results, i.(*Ingredient))
	}

	return results
}
