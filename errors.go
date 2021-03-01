package core

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// ErrInvalidParams http 400 error code
var ErrInvalidParams = errors.New("invalid params")

// ErrObjectConflict http 409 conflict
var ErrObjectConflict = errors.New("object conflict")

// ErrObjectNotFound http 404 not found
var ErrObjectNotFound = errors.New("object not found")

// ErrServerError http 500 server error
var ErrServerError = errors.New("server error")

// ErrRemoteError http 502 server error
var ErrRemoteError = errors.New("remote error")

// ErrNoAuth http 401 no auth
var ErrNoAuth = errors.New("no auth")

// ErrNoPermission http 403 no permission
var ErrNoPermission = errors.New("no permission")

// ErrPaymentRequired 402 Payment Required
var ErrPaymentRequired = errors.New("payment required")

// JSONError return json data struct
func JSONError(code int, err error) ([]byte, error) {
	content, err := jsoniter.Marshal(H{
		"error": H{
			"code":    code,
			"message": err.Error(),
		},
	})

	if err != nil {
		return nil, err
	}

	return content, nil
}
