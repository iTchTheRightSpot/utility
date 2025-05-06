package middleware

import (
	"embed"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/utility/utils"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//go:embed html
var frontend embed.FS

func fileSystem() (http.FileSystem, error) {
	build, err := fs.Sub(frontend, "html")
	if err != nil {
		fmt.Print(err.Error())
		return nil, errors.New("error loading html")
	}
	return http.FS(build), nil
}

func TestMiddleware(t *testing.T) {
	t.Parallel()

	lg := utils.DevLogger("UTC")

	t.Run("SPA middleware", func(t *testing.T) {
		t.Run("should load api route", func(t *testing.T) {
			t.Parallel()

			// given
			fsys, err := fileSystem()
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			m := Middleware{Logger: lg, ApiPrefix: "/api/", Fs: fsys}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/dummy", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/dummy", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.SPA(mux).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusNoContent {
				t.Errorf("expected %d, given %d", http.StatusNoContent, rr.Code)
				t.FailNow()
			}
		})

		t.Run("should load html", func(t *testing.T) {
			t.Parallel()

			// given
			filesys, err := fileSystem()
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			m := Middleware{Logger: lg, ApiPrefix: "/api/", Fs: filesys}

			mux := http.NewServeMux()
			mux.Handle("/", http.StripPrefix("/", http.FileServer(m.Fs)))

			req := httptest.NewRequest(http.MethodGet, "/example.html", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.SPA(mux).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.FailNow()
			}
		})

		t.Run("should load single page", func(t *testing.T) {
			t.Parallel()

			// given
			filesys, err := fileSystem()
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			m := Middleware{Logger: lg, ApiPrefix: "/api/", Fs: filesys}

			mux := http.NewServeMux()
			mux.Handle("/", http.StripPrefix("/", http.FileServer(m.Fs)))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.SPA(mux).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.FailNow()
			}
		})
	})

	t.Run("error middleware", func(t *testing.T) {
		t.Parallel()

		t.Run("route not found", func(t *testing.T) {
			t.Parallel()

			// given
			m := Middleware{Logger: lg}
			mux := http.NewServeMux()

			req := httptest.NewRequest(http.MethodPost, "/path", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.Error(mux).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
				t.FailNow()
			}

			str := strings.TrimSpace(rr.Body.String())
			s := `{"message":"route not found"}`
			if str != s {
				t.Errorf("expect equal, expect %s, got %s", s, str)
				t.FailNow()
			}
		})

		t.Run("method not found", func(t *testing.T) {
			t.Parallel()

			// given
			m := Middleware{Logger: lg}
			mux := http.NewServeMux()

			mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.Error(mux).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, rr.Code)
				t.FailNow()
			}

			str := strings.TrimSpace(rr.Body.String())
			s := `{"message":"method not allowed"}`
			if str != s {
				t.Errorf("expect equal, expect %s, got %s", s, str)
				t.FailNow()
			}
		})

		t.Run("request timeout", func(t *testing.T) {
			t.Parallel()

			// given
			m := Middleware{Logger: lg}
			mux := http.NewServeMux()
			han := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2)
				w.WriteHeader(200)
			})
			mux.Handle("GET /path", http.TimeoutHandler(han, 1, "should time out"))

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			rr := httptest.NewRecorder()

			// method to test
			m.Log(m.Error(mux)).ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusServiceUnavailable {
				t.Errorf("expected status code %d, got %d", http.StatusServiceUnavailable, rr.Code)
				t.FailNow()
			}

			get := rr.Header().Get("Content-Type")
			fmt.Print(get)
			str := strings.TrimSpace(rr.Body.String())
			s := `{"message":"request timeout"}`
			if str != s {
				t.Errorf("expect not equal, expect %s, got %s", s, str)
				t.FailNow()
			}
		})
	})

	t.Run("should recover from panic", func(t *testing.T) {
		t.Parallel()

		// given
		m := Middleware{Logger: lg}
		mux := http.NewServeMux()
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			panic("simulate error")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		// method to test
		m.Panic(mux).ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
			t.FailNow()
		}

		str := strings.TrimSpace(rr.Body.String())
		s := `{"message":"server error"}`
		if str != s {
			t.Errorf("expect equal, expect %s, got %s", s, str)
			t.FailNow()
		}
	})
}