package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

// Minimal info we need to return history
type SnapshotMeta struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
}

func GetSnapshotsByHost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Extract host IP from path (/hosts/{ip}/snapshots) <- how our route looks
		path := strings.TrimPrefix(r.URL.Path, "/hosts/")
		parts := strings.SplitN(path, "/", 2)
		//URL verification
		if len(parts) < 2 || parts[1] != "snapshots" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		hostIP := parts[0]
		// Query DB
		rows, err := db.Query(
			"SELECT id, timestamp FROM snapshots WHERE host_ip = ? ORDER BY timestamp ASC",
			hostIP,
		)
		if err != nil {
			http.Error(w, "db query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var history []SnapshotMeta
		for rows.Next() {
			var meta SnapshotMeta
			if err := rows.Scan(&meta.ID, &meta.Timestamp); err != nil {
				http.Error(w, "db scan failed", http.StatusInternalServerError)
				return
			}
			history = append(history, meta)
		}
		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(history); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
