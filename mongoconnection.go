package fridgedoordatabase

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection represents a connection to a database
type Connection interface {
	Collection(database string, collection string) *mongo.Collection
	Disconnect() error
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

func (db *mongoConnection) Collection(database string, collection string) *mongo.Collection {
	return db.client.Database(database).Collection(collection)
}

func (db *mongoConnection) Disconnect() error {
	err := db.client.Disconnect(context.TODO())

	return err
}
