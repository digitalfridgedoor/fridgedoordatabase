package userview

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestUpdateTags(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")

	duration3s, _ := time.ParseDuration("10s")
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration3s)
	defer cancelFunc()
	connect := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connect)

	username := "TestUser"

	view, err := Create(context.Background(), username)

	assert.Nil(t, err)

	assert.Equal(t, username, view.Username)
	assert.NotNil(t, view.Collections)

	viewID := view.ID
	tag := "tag"
	err = AddTag(context.Background(), view.ID.Hex(), tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 1, len(view.Tags))
	assert.Equal(t, tag, view.Tags[0])

	err = AddTag(context.Background(), view.ID.Hex(), tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 1, len(view.Tags))
	assert.Equal(t, tag, view.Tags[0])

	err = RemoveTag(context.Background(), viewID.Hex(), tag)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)
	assert.Equal(t, 0, len(view.Tags))

	err = Delete(ctx, username)
	assert.Nil(t, err)

	_, err = GetByUsername(context.Background(), username)
	assert.NotNil(t, err)
}
