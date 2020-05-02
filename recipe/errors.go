package recipe

import "errors"

var errNotConnected = errors.New("Not connected")
var errDuplicate = errors.New("Duplicate")
var errSubRecipe = errors.New("Invalid subrecipe reference")
var errUnauthorised = errors.New("Cannot perform requested action on this recipe")
