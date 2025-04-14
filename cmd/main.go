package main

import (
	"context"
	"github.com/iTchTheRightSpot/utility/utils"
	"net/http"
)

func main() {
	lg := utils.DevLogger("America/Toronto")
	m := utils.Middleware{Logger: lg}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lg.Log(r.Context(), "base path hit")
		w.WriteHeader(200)
	})

	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		lg.Log(r.Context(), "api path hit")
		w.WriteHeader(200)
	})

	server := http.Server{Addr: ":8080", Handler: m.Initialize(mux)}

	lg.Log(context.Background(), "server listening on Post 8080")
	lg.Fatal(server.ListenAndServe())
}