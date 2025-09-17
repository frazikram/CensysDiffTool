package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hostdiff/api/internal/models"
	"hostdiff/api/internal/services"
)

//Get Diff, comparing 2 snapshots by ID

func GetDiff(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Get IDs from query params
		q := r.URL.Query()
		id1 := q.Get("snap1")
		id2 := q.Get("snap2")
		if id1 == "" || id2 == "" {
			http.Error(w, "missing snap1 or snap2 query param", http.StatusBadRequest)
			return
		}
		//Get Snapshots
		snap1, err := fetchSnapshot(db, id1)
		if err != nil {
			http.Error(w, "snapshot 1 not found", http.StatusNotFound)
			return
		}
		snap2, err := fetchSnapshot(db, id2)
		if err != nil {
			http.Error(w, "snapshot 2 not found", http.StatusNotFound)
			return
		}
		// Ensure both snapshots belong to the same host
		if snap1.IP != snap2.IP {
			http.Error(w, "snapshots belong to different hosts", http.StatusBadRequest)
			return
		}
		//Compute diff
		diff := services.DiffSnapshots(snap1, snap2)

		//Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(diff)

	}
}

// Helper Function to get and unpack JSON from DB
func fetchSnapshot(db *sql.DB, id string) (models.Snapshot, error) {
	var raw string
	row := db.QueryRow("SELECT data FROM snapshots WHERE id = ?", id)
	err := row.Scan(&raw)
	if err != nil {
		return models.Snapshot{}, err
	}

	var snap models.Snapshot
	if err := json.Unmarshal([]byte(raw), &snap); err != nil {
		return models.Snapshot{}, err
	}

	return snap, nil
}
