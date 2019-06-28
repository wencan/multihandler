package multihandler

/*
 * wrap the http response structure
 * to get the http response data.
 * modify from github.com/gorilla/handlers
 *
 * wencan
 * 2019-06-27
 */

import (
	"bufio"
	"net"
	"net/http"
)

// response wrap http.ResponseWriter
type response struct {
	w      http.ResponseWriter
	status int
	size   int

	wroteHeader bool
}

func (resp *response) Header() http.Header {
	return resp.w.Header()
}

func (resp *response) Write(b []byte) (int, error) {
	if !resp.wroteHeader {
		resp.wroteHeader = true
		resp.status = http.StatusOK
	}

	size, err := resp.w.Write(b)
	resp.size += size
	return size, err
}

func (resp *response) WriteHeader(status int) {
	resp.wroteHeader = true
	resp.w.WriteHeader(status)
	resp.status = status
}

func (resp *response) Status() int {
	return resp.status
}

func (resp *response) Size() int {
	return resp.size
}

// Flush implement http.Flusher
func (resp *response) Flush() {
	f, ok := resp.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// hijacker wrap http.Hijacker
type hijacker struct {
	response
}

func (jacker *hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := jacker.response.w.(http.Hijacker)
	conn, rw, err := h.Hijack()
	if err == nil && jacker.response.status == 0 {
		jacker.response.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}
