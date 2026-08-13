package httpapi

import (
	"errors"
	"math"
	"net/http"
	"strings"
	"testing"
)

func TestWriteJSONReturnsEncodingError(t *testing.T) {
	writer := &failingResponseWriter{}
	err := writeJSON(writer, http.StatusOK, resultResponse{Result: math.NaN()})
	if err == nil || !strings.Contains(err.Error(), "encode JSON response") {
		t.Errorf("writeJSON() error = %v, want encoding error", err)
	}
	if writer.status != 0 {
		t.Errorf("status = %d, want headers to remain uncommitted", writer.status)
	}
}

func TestWriteJSONReturnsWriteError(t *testing.T) {
	writer := &failingResponseWriter{writeError: errors.New("connection closed")}
	err := writeJSON(writer, http.StatusOK, resultResponse{Result: 42})
	if err == nil || !strings.Contains(err.Error(), "write JSON response") {
		t.Errorf("writeJSON() error = %v, want write error", err)
	}
}

type failingResponseWriter struct {
	header     http.Header
	status     int
	writeError error
}

func (w *failingResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *failingResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, w.writeError
}
