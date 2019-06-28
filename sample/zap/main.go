package main

import (
	"log"
	"net/http"

	"github.com/wencan/multihandler"
	zaplog "github.com/wencan/multihandler/zap"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Println(err)
		return
	}

	handle := func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/panic" {
			panic("BOOM!")
		}
		w.Write([]byte("Hello world"))
	}

	handler := multihandler.NewMultiHandler(zaplog.NewZapLogging(logger), http.HandlerFunc(handle))

	err = http.ListenAndServe(":8080", handler)
	logger.Error("service stop", zap.Error(err))
}
