package error

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrRecordAlreadyExists = errors.New("record already exists")
var ErrRecordNotFound = errors.New("record not found")
var ErrInvalidEmailPassword = errors.New("invalid email/password")
var ErrEmailNotVerified = errors.New("email not verified")
var ErrNotAuthorize = errors.New("not authorize to perform this action")
