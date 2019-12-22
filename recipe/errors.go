package recipe

import "errors"

var errNotConnected = errors.New("Not connected")
var errDuplicate = errors.New("Duplicate")
var errSubRecipe = errors.New("Invalid subrecipe reference")
