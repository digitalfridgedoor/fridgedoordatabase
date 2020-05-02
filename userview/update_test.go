package userview

import (
	"context"
	"testing"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"github.com/digitalfridgedoor/fridgedoordatabase/dfdtesting"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateTags(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetUserViewFindPredicate(func(uv *dfdmodels.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})

	username := "TestUser"

	view, err := Create(context.Background(), username)

	assert.Nil(t, err)

	assert.Equal(t, username, view.Username)

	viewID := view.ID
	tag := "tag"
	err = AddTag(context.Background(), view.ID, tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 1, len(view.Tags))
	assert.Equal(t, tag, view.Tags[0])

	err = AddTag(context.Background(), view.ID, tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 1, len(view.Tags))
	assert.Equal(t, tag, view.Tags[0])

	err = RemoveTag(context.Background(), viewID, tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 0, len(view.Tags))

	err = Delete(context.TODO(), username)
	assert.Nil(t, err)

	_, err = GetByUsername(context.Background(), username)
	assert.NotNil(t, err)
}
