package ingredient

import (
	"context"
	"log"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection can find and parse Ingredients from mongodb
type Connection struct {
	db fridgedoordatabase.Connection
}

// Find does not find one
func (conn *Connection) Find(ctx context.Context) ([]*Ingredient, error) {

	collection := conn.db.Collection("recipeapi", "ingredients")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
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
func (conn *Connection) FindOne(ctx context.Context, id string) (*Ingredient, error) {

	collection := conn.db.Collection("recipeapi", "ingredients")

	// Pass these options to the Find method
	findOneOptions := options.FindOne()
	// findOneOptions.SetLimit(2)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	singleResult := collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, findOneOptions)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Ingredient{})

	return ing.(*Ingredient), err
}

// IngredientByParentID returns an array of ingredients with the given parentID
func (conn *Connection) IngredientByParentID(ctx context.Context, parentID primitive.ObjectID) []*Ingredient {

	collection := conn.db.Collection("recipeapi", "ingredients")

	cur, err := collection.Find(ctx, bson.M{"parentId": parentID})
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
