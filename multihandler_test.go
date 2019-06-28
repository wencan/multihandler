package multihandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/wencan/multihandler/mock_multihandler"
)

func makeWriteLogSubtest(url string, handlerFunc http.HandlerFunc) func(t *testing.T) {
	return func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loggingEngine := mock_multihandler.NewMockLoggingEngine(ctrl)

		var request *http.Request
		var respStatus, respBodyBytesSent int
		loggingEngine.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(req *http.Request, status, bodyBytesSent int, timestamp time.Time) {
			request = req
			respStatus = status
			respBodyBytesSent = bodyBytesSent
		}).AnyTimes()

		mockReq := httptest.NewRequest(http.MethodGet, url, nil)
		recorder := httptest.NewRecorder()
		handler := NewMultiHandler(loggingEngine, http.HandlerFunc(handlerFunc))
		handler.ServeHTTP(recorder, mockReq)

		if respStatus != recorder.Code {
			t.Errorf("http response status not equal, got: %d, expect:%d", respStatus, recorder.Code)
		}
		if respBodyBytesSent != recorder.Body.Len() {
			t.Errorf("http body bytes sent length not equal, got: %d, expect:%d", respBodyBytesSent, recorder.Body.Len())
		}
		if request != mockReq {
			t.Errorf("http request not equal, got: %p, expect:%p", request, mockReq)
		}
	}
}

func TestWriteLog(t *testing.T) {
	okhandle := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}
	t.Run("ok", makeWriteLogSubtest("/ok", okhandle))

	badhandle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad"))
	}
	t.Run("bad", makeWriteLogSubtest("/bad", badhandle))
}

func TestWritePanicLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggingEngine := mock_multihandler.NewMockLoggingEngine(ctrl)

	var request *http.Request
	var respStatus int
	var recovered interface{}
	loggingEngine.EXPECT().WritePanic(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(req *http.Request, status, bodyBytesSent int, timestamp time.Time, r interface{}) {
		request = req
		respStatus = status
		recovered = r
	}).AnyTimes()

	panicError := "this is a test"
	panicHandle := func(w http.ResponseWriter, r *http.Request) {
		panic(panicError)
	}

	mockReq := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()
	handler := NewMultiHandler(loggingEngine, http.HandlerFunc(panicHandle))
	handler.ServeHTTP(recorder, mockReq)

	if respStatus != recorder.Code {
		t.Errorf("http response status not equal, got: %d, expect:%d", respStatus, recorder.Code)
	}
	if request != mockReq {
		t.Errorf("http request not equal, got: %p, expect:%p", request, mockReq)
	}
	if recovered != panicError {
		t.Errorf("error not equal, got: %v, expect:%v", recovered, panicError)
	}
}
