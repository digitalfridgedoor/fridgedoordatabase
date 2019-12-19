package fridgedoordatabase

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// Collection wraps a connected mongo collection
type Collection struct {
	MongoCollection *mongo.Collection
}

// Connect connects to mongo
func Connect(ctx context.Context, connectionString string) bool {
	if mongoClient != nil {
		// assume connected
		return true
	}

	return connect(ctx, connectionString)
}

func connect(ctx context.Context, connectionString string) bool {
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

	mongoClient = client

	return true
}

// CreateCollection gets a wrapped reference to a mongo collection
func CreateCollection(database string, collection string) (bool, *Collection) {
	if mongoClient == nil {
		return false, nil
	}
	return true, &Collection{mongoClient.Database(database).Collection(collection)}
}

// Disconnect removes the current connection to mongo
func Disconnect() error {
	if mongoClient == nil {
		return nil
	}
	err := mongoClient.Disconnect(context.TODO())
	mongoClient = nil

	return err
}
