package main

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
		// the header can contain multiple IPs, so take the first one
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
}

func (dep *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := &RequestBody{
			Id:     uuid.NewString(),
			Ip:     dep.clientIP(r),
			Method: r.Method,
			Path:   r.Method,
		}
		r = r.WithContext(context.WithValue(r.Context(), RequestKey, b))
		dep.Logger.Log(r.Context(), "Request")
		obj := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(obj, r)
		dep.Logger.Log(r.Context(), fmt.Sprintf("Response Status: %d", obj.statusCode))
	})
}
