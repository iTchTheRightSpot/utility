package log

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type Middleware struct {
	Logger ILogger
}

func (dep *Middleware) Initialize(router *http.ServeMux) http.Handler {
	return dep.Log(router)
}

// https://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request
func (dep *Middleware) clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
}

func (dep *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := dep.Logger.Date()
		b := &RequestBody{
			Id:     uuid.NewString(),
			Ip:     dep.clientIP(r),
			Method: r.Method,
			Path:   r.URL.Path,
		}
		r = r.WithContext(context.WithValue(r.Context(), RequestKey, b))
		obj := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(obj, r)
		str := fmt.Sprintf("Response Status: %d | Duration: %v second(s)", obj.statusCode, dep.Logger.Date().Sub(start).Seconds())
		dep.Logger.Log(r.Context(), str)
	})
}
