package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"hostdiff/api/internal/models"
)

// function handles uploading a new snapshot
func PostSnapshot(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow POST
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

		// Parse body into Snapshot struct
		var snapshot models.Snapshot
		if err := json.Unmarshal(body, &snapshot); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Basic validation
		if snapshot.IP == "" || snapshot.Timestamp == "" {
			http.Error(w, "missing ip or timestamp", http.StatusBadRequest)
			return
		}

		// Insert into DB (store JSON blob as string)
		query := `INSERT INTO snapshots (host_ip, timestamp, data) VALUES (?, ?, ?)`
		res, err := db.Exec(query, snapshot.IP, snapshot.Timestamp, string(body))
		if err != nil {
			// Catch duplicate snapshot (UNIQUE constraint)
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Snapshot already exists for this host and timestamp", http.StatusConflict)
				return
			}

			// Other DB errors
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		id, _ := res.LastInsertId()

		// Return success response
		resp := map[string]any{
			"id":     id,
			"status": "stored",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
