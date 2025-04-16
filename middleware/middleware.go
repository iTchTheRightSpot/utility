package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"strings"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type Middleware struct {
	Logger    utils.ILogger
	Fs        http.FileSystem
	ApiPrefix string
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
		b := &utils.RequestBody{
			Id:     uuid.NewString(),
			Ip:     dep.clientIP(r),
			Method: r.Method,
			Path:   r.URL.Path,
		}
		r = r.WithContext(context.WithValue(r.Context(), utils.RequestKey, b))
		obj := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(obj, r)
		str := fmt.Sprintf("Response Status: %d | Duration: %v second(s)", obj.statusCode, dep.Logger.Date().Sub(start).Seconds())
		dep.Logger.Log(r.Context(), str)
	})
}

// SPA loads single page api
func (dep *Middleware) SPA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, dep.ApiPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		file, err := dep.Fs.Open(r.URL.Path)
		if err == nil {
			defer func(file http.File) {
				if err = file.Close(); err != nil {
					dep.Logger.Critical(r.Context(), err.Error())
				}
			}(file)

			if stat, statErr := file.Stat(); statErr == nil && !stat.IsDir() {
				http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
				return
			}
		}

		// If file doesn't exist, fallback to index.html for SPA
		indexFile, err := dep.Fs.Open("index.html")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer func(indexFile http.File) {
			if err = indexFile.Close(); err != nil {
				dep.Logger.Critical(r.Context(), err.Error())
			}
		}(indexFile)

		if stat, statErr := indexFile.Stat(); statErr == nil {
			http.ServeContent(w, r, stat.Name(), stat.ModTime(), indexFile)
			return
		}

		next.ServeHTTP(w, r)
	})
}