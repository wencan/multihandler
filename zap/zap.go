package zap

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ZapLogging Implement of LoggingEngine
type ZapLogging struct {
	logger *zap.Logger
}

// NewZapLogging create ZapLogging object
func NewZapLogging(logger *zap.Logger) *ZapLogging {
	return &ZapLogging{
		logger: logger,
	}
}

func buildZapFields(req *http.Request, status, bodyBytesSent int, timestamp time.Time) []zap.Field {
	fields := []zap.Field{}
	fields = append(fields, zap.String("remote_addr", req.RemoteAddr))
	user, _, ok := req.BasicAuth()
	if ok {
		fields = append(fields, zap.String("remote_user", user))
	}
	fields = append(fields, zap.Time("time_local", timestamp))
	fields = append(fields, zap.String("request_method", req.Method))
	fields = append(fields, zap.String("request_uri", req.RequestURI))
	fields = append(fields, zap.String("server_protocol", req.Proto))
	fields = append(fields, zap.Int("status", status))
	fields = append(fields, zap.Int("body_bytes_sent", bodyBytesSent))
	if req.Referer() != "" {
		fields = append(fields, zap.String("http_referer", req.Referer()))
	}
	if req.UserAgent() != "" {
		fields = append(fields, zap.String("http_user_agent", req.UserAgent()))
	}
	fields = append(fields, zap.String("elapsed_time", time.Now().Sub(timestamp).String()))
	return fields
}

// Write logs a message.
func (logging *ZapLogging) Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error {
	fields := buildZapFields(req, status, bodyBytesSent, timestamp)

	statusText := http.StatusText(status)
	switch status / 100 {
	case 1:
		logging.logger.Info(statusText, fields...)
	case 2:
		logging.logger.Info(statusText, fields...)
	case 3:
		logging.logger.Info(statusText, fields...)
	case 4:
		logging.logger.Warn(statusText, fields...)
	case 5:
		logging.logger.Error(statusText, fields...)
	default:
		logging.logger.Warn("Unknown HTTP status", fields...)
	}
	return nil
}

// WritePanic logs a message at DPanicLevel. If the logger is in development mode, it then panics.
func (logging *ZapLogging) WritePanic(req *http.Request, status, bodyBytesSent int, timestamp time.Time, recovered interface{}) error {
	fields := buildZapFields(req, status, bodyBytesSent, timestamp)
	fields = append(fields, zap.Any("recovered", recovered))
	logging.logger.DPanic("recover a panic", fields...)
	return nil
}
