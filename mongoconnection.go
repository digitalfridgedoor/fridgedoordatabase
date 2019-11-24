package fridgedoordatabase

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ingredient represents a node in the ingredient tree
type Ingredient struct {
	ID       string             `json:"_id"`
	Name     string             `json:"name"`
	AddedOn  time.Time          `json:"addedOn"`
	ParentID primitive.ObjectID `json:"parentId"`
}

// Connection represents a connection to a database
type Connection interface {
	Find(ctx context.Context) []*Ingredient
	Disconnect() error
	IngredientByParentID(ctx context.Context, parentID primitive.ObjectID) []*Ingredient
}

type mongoConnection struct {
	client *mongo.Client
}

// Connect connects to mongo
func Connect(ctx context.Context, connectionString string) Connection {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return &mongoConnection{client}
}

func (db *mongoConnection) Disconnect() error {
	err := db.client.Disconnect(context.TODO())

	return err
}

func (db *mongoConnection) FindOne(ctx context.Context) Ingredient {
	collection := db.client.Database("recipeapi").Collection("ingredients")

	result := collection.FindOne(ctx, bson.D{{}})

	var ing Ingredient

	err := result.Decode(&ing)
	if err != nil {
		log.Fatal(err)
	}

	return ing
}

// IngredientByParentID returns an array of ingredients with the given parentID
func (db *mongoConnection) IngredientByParentID(ctx context.Context, parentID primitive.ObjectID) []*Ingredient {
	collection := db.client.Database("recipeapi").Collection("ingredients")

	cur, err := collection.Find(ctx, bson.M{"parentId": parentID})
	if err != nil {
		log.Fatal(err)
	}

	return parseIngredients(ctx, cur)
}

func (db *mongoConnection) Find(ctx context.Context) []*Ingredient {

	collection := db.client.Database("recipeapi").Collection("ingredients")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	return parseIngredients(ctx, cur)
}

func parseIngredients(ctx context.Context, cur *mongo.Cursor) []*Ingredient {
	var results []*Ingredient

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem Ingredient
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(ctx)

	return results
}
