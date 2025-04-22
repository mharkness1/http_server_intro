package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM != "dev" {
		w.WriteHeader(403)
	}
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(w, "Counter reset to 0")
	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
