package main

import (
	"fmt"
	"hostdiff/api/internal/handlers"
	"hostdiff/api/internal/storage"
	"log"
	"net/http"
)

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
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
