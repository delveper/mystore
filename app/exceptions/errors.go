package exceptions

import "errors"

var ErrInvalidData = errors.New("invalid data")
var ErrEmptyContext = errors.New("context is empty")
var ErrNotFound = errors.New("not found")
var ErrUnexpected = errors.New("unexpected error")
var ErrDeadline = errors.New("deadline exceeded")
var ErrRecordExists = errors.New("record already exists")

var ErrMerchantNotFound = errors.New("merchant not found")
