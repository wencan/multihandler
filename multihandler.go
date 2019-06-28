package multihandler

/*
 * HTTP middleware that supports multiple functions.
 *
 * wencan
 * 2019-06-27
 */

import (
	"net/http"
	"time"
)

// MultiHandler HTTP middleware that supports multiple functions.
type MultiHandler struct {
	loggingEngine LoggingEngine

	next http.Handler
}

// NewMultiHandler Create MultiHandler object.
func NewMultiHandler(loggingEngine LoggingEngine, next http.Handler) *MultiHandler {
	return &MultiHandler{
		loggingEngine: loggingEngine,
		next:          next,
	}
}

func (handler *MultiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// wrap response
	resp := &response{w: w}

	timestamp := time.Now()
	var recoverd interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				recoverd = r

				// write http status 500
				if resp.Status() == 0 {
					resp.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		// next
		handler.next.ServeHTTP(resp, req)
	}()

	// log
	if recoverd != nil {
		handler.loggingEngine.WritePanic(req, resp.Status(), resp.Size(), timestamp, recoverd)
	} else {
		handler.loggingEngine.Write(req, resp.Status(), resp.Size(), timestamp)
	}
}
