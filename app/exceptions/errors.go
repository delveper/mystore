package exceptions

import "errors"

var ErrInvalidData = errors.New("invalid data, expected nil")
var ErrNotAuthorized = errors.New("not authorized")

var ErrEmptyContext = errors.New("context is empty")
var ErrRecordNotFound = errors.New("record not found")
var ErrUnexpected = errors.New("unexpected error")
var ErrDeadline = errors.New("deadline exceeded")
var ErrRecordExists = errors.New("record already exists")

/*
var ErrDuplicateID = errors.New("id already exists")
var ErrDuplicateEmail = errors.New("email is already taken")
var ErrDuplicatePhone = errors.New("phone already exists")
*/
