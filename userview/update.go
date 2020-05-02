package userview

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetNickname updates the users nickname
func SetNickname(ctx context.Context, view *dfdmodels.UserView, nickname string) error {

	if nickname == "" || view.Nickname == nickname {
		return nil
	}

	view.Nickname = nickname

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	return coll.c.UpdateByID(ctx, &view.ID, view)
}

// AddTag adds a tag to users list if it isn't already there
func AddTag(ctx context.Context, id *primitive.ObjectID, tag string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	view, err := coll.findOne(ctx, id)
	if err != nil {
		return err
	}

	hasTag := false
	for _, t := range view.Tags {
		if t == tag {
			hasTag = true
		}
	}

	if !hasTag {
		view.Tags = append(view.Tags, tag)
	}

	return coll.c.UpdateByID(ctx, id, view)
}

// RemoveTag removes a tag from a users list
func RemoveTag(ctx context.Context, id *primitive.ObjectID, tag string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	view, err := coll.findOne(ctx, id)
	if err != nil {
		return err
	}

	view.Tags = filterTags(view.Tags, tag)

	return coll.c.UpdateByID(ctx, id, view)
}

func filterTags(tags []string, tagToRemove string) []string {
	filtered := []string{}

	for _, tag := range tags {
		if tag != tagToRemove {
			filtered = append(filtered, tag)
		}
	}

	return filtered
}
