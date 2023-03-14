package exceptions

import "errors"

var ErrInvalidData = errors.New("invalid data, expected nil")
var ErrPermissions = errors.New("error permissions")
var ErrInvalidQuery = errors.New("invalid query")
var ErrNotAuthorized = errors.New("not authorized")
