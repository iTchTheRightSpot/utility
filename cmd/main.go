package main

import (
	"context"
	"encoding/json"
	"github.com/iTchTheRightSpot/utility/middleware"
	"github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"time"
)

func main() {
	lg := utils.DevLogger("UTC")
	m := middleware.Middleware{Logger: lg, ApiPrefix: "/api"}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(200)
		if _, err := w.Write([]byte("slow response")); err != nil {
			if err = json.NewEncoder(w).Encode(&utils.ServerError{}); err != nil {
				utils.ErrorResponse(w, err)
			}
		}
	})

	srv := http.Server{
		Addr:              ":8080",
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       3 * time.Second,
		Handler:           m.Log(m.Timeout(3*time.Second, mux)),
	}

	lg.Log(context.Background(), "server listening on port 8080")
	lg.Fatal(srv.ListenAndServe())
}
