package storages

import "errors"

var ErrJamNotFound = errors.New("not found jam")
var ErrJamURLConflit = errors.New("jam with this url already exists")
var ErrGameNotFound = errors.New("not found game")
