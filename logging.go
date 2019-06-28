package multihandler

/*
 * http logging handler
 *
 * wencan
 * 2019-06-27
 */

import (
	"net/http"
	"time"
)

// LoggingEngine interface of logger wrap
type LoggingEngine interface {
	Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error

	WritePanic(req *http.Request, status, bodyBytesSent int, timestamp time.Time, recovered interface{}) error
}
