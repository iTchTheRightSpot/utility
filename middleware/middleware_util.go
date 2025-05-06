package middleware

import (
	"net/http"
)

type logWriter struct {
	http.ResponseWriter
	code int
}

func (w *logWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

type errorWriter struct {
	http.ResponseWriter
	code     int
	override bool
}

func (w *errorWriter) WriteHeader(code int) {
	// keep track of issue in case there is an easier way to
	// update this behaviour https://github.com/golang/go/issues/65648
	t := w.Header().Get("Content-Type")
	if t == "text/plain; charset=utf-8" || (t == "" && code == http.StatusServiceUnavailable) {
		w.Header().Set("Content-Type", "application/json")
		w.override = true
	}
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *errorWriter) Write(body []byte) (int, error) {
	if w.override {
		switch w.code {
		case http.StatusNotFound:
			body = []byte(`{"message":"route not found"}`)
		case http.StatusMethodNotAllowed:
			body = []byte(`{"message":"method not allowed"}`)
		case http.StatusServiceUnavailable:
			body = []byte(`{"message":"request timeout"}`)
		}
	}
	return w.ResponseWriter.Write(body)
}