package storages

import "errors"

var ErrJamURLConflit = errors.New("jam with this url already exists")
var ErrNotFound = errors.New("not found")
