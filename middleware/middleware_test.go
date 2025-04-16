package middleware

import (
	"embed"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/utility/utils"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
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
}