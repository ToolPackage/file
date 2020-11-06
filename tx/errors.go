package tx

import "github.com/pkg/errors"

func InvalidPacketError(any interface{}) error {
	var msg string
	if err, ok := any.(error); ok {
		msg = err.Error()
	} else if str, ok := any.(string); ok {
		msg = str
	} else {
		return errors.New("failed to get root cause")
	}
	return errors.New("invalid packet error: " + msg)
}
