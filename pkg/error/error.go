package error

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrRecordNotFound = errors.New("record not found")
var ErrInvalidEmailPassword = errors.New("invalid email/password")
var ErrEmailNotVerified = errors.New("email not verified")
