package recipe

import (
	"context"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/user"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new recipe with given name
func (conn *Connection) Create(ctx context.Context, userID string, name string) (*Recipe, error) {

	u := user.New(conn.db)
	userInfo, err := u.FindOne(ctx, userID)
	if err != nil {
		return nil, err
	}

	collection := conn.collection()

	insertOneOptions := options.InsertOne()

	recipe := &Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userInfo.ID,
	}

	_, err = collection.InsertOne(ctx, recipe, insertOneOptions)
	if err != nil {
		return nil, err
	}
	return nil, nil

	// insertOneResult.InsertedID

	// userInfo.Recipes

	// ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})

	// return ing.(*Recipe), err
}
