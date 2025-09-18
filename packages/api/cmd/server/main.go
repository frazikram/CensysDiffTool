package main

import (
	"fmt"
	"hostdiff/api/internal/handlers"
	"hostdiff/api/internal/storage"
	"log"
	"net/http"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend (adjust to your dev server if you want stricter rules)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func main() {
	//Initializing DB
	db := storage.InitDB("snapshots.db")
	defer db.Close()
	mux := http.NewServeMux()

	//Simple health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	// POST /snapshots
	mux.Handle("/snapshots", handlers.PostSnapshot(db))
	// GET History
	mux.Handle("/hosts/", handlers.GetSnapshotsByHost(db))
	//GET Diff
	mux.Handle("/diff", handlers.GetDiff(db))
	//Get All Hosts
	mux.HandleFunc("/all-hosts", handlers.HostsHandler(db))
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
