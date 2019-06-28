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

// NoOpLogging Implement of LoggingEngine. It do nothing.
type NoOpLogging struct {
}

func (logging *NoOpLogging) Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error {
	return nil
}

func (logging *NoOpLogging) WritePanic(req *http.Request, status, bodyBytesSent int, timestamp time.Time, recovered interface{}) error {
	return nil
}
