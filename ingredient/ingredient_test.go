package ingredient

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {

	capital, err := FindByName(context.Background(), "C")
	lowercase, err := FindByName(context.Background(), "c")

	assert.Nil(t, err)
	assert.Equal(t, len(capital), len(lowercase))
}

func TestFindOne(t *testing.T) {

	id, err := primitive.ObjectIDFromHex("5d8f744446106c8aee8cde37")
	assert.Nil(t, err)

	ing, err := FindOne(context.Background(), &id)

	assert.Nil(t, err)
	assert.NotNil(t, ing)
	assert.Equal(t, "5dac764fa0b9423b0090a898", ing.ParentID.Hex())
	assert.Equal(t, "5d8f744446106c8aee8cde37", ing.ID.Hex())
	assert.Equal(t, "Chicken thighs", ing.Name)
}
