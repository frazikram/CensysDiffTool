package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hostdiff/api/internal/models"
)

// function handles uploading a new snapshot
func PostSnapshot(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//only allow post
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		// Parse Body for snapshot struct
		var snapshot models.Snapshot
		if err := json.Unmarshal(body, &snapshot); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		//Basic Validation to check if data is good
		if snapshot.IP == "" || snapshot.Timestamp == "" {
			http.Error(w, "missing ip or timestamp", http.StatusBadRequest)
			return
		}
		//Inserting into db (json blob as string)
		query := `INSERT INTO snapshots (host_ip, timestamp, data) VALUES (?, ?, ?)`
		res, err := db.Exec(query, snapshot.IP, snapshot.Timestamp, string(body))
		if err != nil {
			http.Error(w, fmt.Sprintf("db insert failed: %v", err), http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		// Return success response
		resp := map[string]any{
			"id":     id,
			"status": "stored",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
