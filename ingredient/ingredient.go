package ingredient

import (
	"context"
	"log"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a ingredient-wrapped collection
type Collection struct {
	collection *fridgedoordatabase.Collection
}

// New creates an instance of ingredient.Collection
func New(db fridgedoordatabase.Connection) *Collection {
	return &Collection{db.Collection("recipeapi", "ingredients")}
}

func (coll *Collection) mongoCollection() *mongo.Collection {
	return coll.collection.MongoCollection
}

// Find does not find one
func (coll *Collection) Find(ctx context.Context) ([]*Ingredient, error) {

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	cur, err := coll.mongoCollection().Find(ctx, bson.D{{}}, findOptions)
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
