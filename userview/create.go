package userview

import (
	"context"
	"errors"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
)

var errUserExists = errors.New("User exists")

// Create creates a new userview for a user
func Create(ctx context.Context, username string) (*dfdmodels.UserView, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	_, err := coll.getByUsername(ctx, username)
	if err == nil {
		// found user with that username
		return nil, errUserExists
	}

	view := &dfdmodels.UserView{
		Username: username,
	}

	v, err := coll.c.InsertOneAndFind(ctx, view, &dfdmodels.UserView{})
	if err != nil {
		return nil, err
	}

	return v.(*dfdmodels.UserView), nil
}

// Delete removes a userview for a user
func Delete(ctx context.Context, username string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	view, err := coll.getByUsername(ctx, username)
	if err != nil {
		return err
	}

	return coll.c.DeleteByID(ctx, &view.ID)
}
