package servers

import "errors"

var ErrPasswordsNotMatch = errors.New("passwords not match")
var ErrFailedAuth = errors.New("incorrect username or password")
