package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// HostsHandler returns a handler that lists unique hosts
func HostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT DISTINCT host_ip FROM snapshots ORDER BY host_ip")
		if err != nil {
			http.Error(w, "failed to fetch hosts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var hosts []string
		for rows.Next() {
			var ip string
			if err := rows.Scan(&ip); err == nil {
				hosts = append(hosts, ip)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hosts)
	}
}
