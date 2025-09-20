package cutil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

type WriteCounter struct {
	http.ResponseWriter
	count      int64
	started    time.Time
	statusCode int
}

func NewWriteCounter(w http.ResponseWriter) *WriteCounter {
	return &WriteCounter{
		ResponseWriter: w,
		started:        util.TimeCurrent(),
	}
}

func (w *WriteCounter) Write(buf []byte) (int, error) {
	n, err := w.ResponseWriter.Write(buf)
	atomic.AddInt64(&w.count, int64(n))
	return n, err
}

func (w *WriteCounter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *WriteCounter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.Header().Set("X-Runtime", fmt.Sprintf("%.6f", time.Since(w.started).Seconds()))
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *WriteCounter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, err := util.Cast[http.Hijacker](w.ResponseWriter)
	if err != nil {
		return nil, nil, errors.Errorf("can't process response of type [%T]", w.ResponseWriter)
	}
	return h.Hijack()
}

func (w *WriteCounter) Count() int64 {
	return atomic.LoadInt64(&w.count)
}

func (w *WriteCounter) Started() time.Time {
	return w.started
}

func (w *WriteCounter) StatusCode() int {
	return w.statusCode
}

func (w *WriteCounter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
