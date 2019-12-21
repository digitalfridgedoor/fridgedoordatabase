package recipe

import "errors"

var errNotConnected = errors.New("Not connected")
var errDuplicate = errors.New("Duplicate")
var errSubRecipe = errors.New("SubRecipes can only have depth 1")
