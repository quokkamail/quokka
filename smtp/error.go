package smtp

import "errors"

var (
	ErrServerClosed = errors.New("smtp: Server closed")
)
