package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"runtime"
	"strings"
)

type wrappedWriter struct {
	http.ResponseWriter
	code             int
	override         bool
	Disable404And405 bool
}

func (w *wrappedWriter) WriteHeader(code int) {
	// keep track of issue in case there is an easier way to
	// update this behaviour https://github.com/golang/go/issues/65648
	if !w.Disable404And405 && w.Header().Get("Content-Type") == "text/plain; charset=utf-8" {
		w.Header().Set("Content-Type", "application/json")
		w.override = true
	}
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrappedWriter) Write(body []byte) (int, error) {
	if w.override {
		switch w.code {
		case http.StatusNotFound:
			body = []byte(`{"message": "route not found"}`)
		case http.StatusMethodNotAllowed:
			body = []byte(`{"message": "method not allowed"}`)
		}
	}
	return w.ResponseWriter.Write(body)
}

type Middleware struct {
	Logger    utils.ILogger
	Fs        http.FileSystem
	ApiPrefix string
	// If true, error response body is overridden for routes or routes methods that do not exist.
	Disable404And405 bool
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
		obj := &wrappedWriter{ResponseWriter: w, code: http.StatusOK, Disable404And405: dep.Disable404And405}
		next.ServeHTTP(obj, r)
		str := fmt.Sprintf("Response Status: %d | Duration: %v second(s)", obj.code, dep.Logger.Date().Sub(start).Seconds())
		dep.Logger.Log(r.Context(), str)
	})
}

// SPA loads single page or frontend pages registered in FileSystem
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

// Panic middleware only handles panic in the main go routine not goroutines spunned within main goroutine
func (dep *Middleware) Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, true)
				buf = buf[:n]
				dep.Logger.Critical(r.Context(), err, string(buf))
				utils.ErrorResponse(w, &utils.ServerError{})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
