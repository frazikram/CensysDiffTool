package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hostdiff/api/internal/storage"
)

//
// --- Helpers ---
//

func setupTestDB(t *testing.T) *sql.DB {
	db := storage.InitDB(":memory:") // in-memory SQLite
	return db
}

func insertSnapshot(t *testing.T, db *sql.DB, ip, ts, data string) int64 {
	res, err := db.Exec(`INSERT INTO snapshots (host_ip, timestamp, data) VALUES (?, ?, ?)`, ip, ts, data)
	if err != nil {
		t.Fatalf("failed to insert snapshot: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

//
// --- Tests ---
//

// POST /snapshots
func TestPostSnapshot(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := PostSnapshot(db)

	body := `{
		"timestamp": "2025-09-15T08:49:45Z",
		"ip": "125.199.235.74",
		"services": []
	}`

	req := httptest.NewRequest(http.MethodPost, "/snapshots", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["status"] != "stored" {
		t.Errorf("expected status=stored, got %v", resp["status"])
	}
}

// GET /hosts/{ip}/snapshots
func TestGetSnapshotsByHost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertSnapshot(t, db, "125.199.235.74", "2025-09-15T08:49:45Z", `{"dummy":"data1"}`)
	insertSnapshot(t, db, "125.199.235.74", "2025-09-16T08:49:45Z", `{"dummy":"data2"}`)

	handler := GetSnapshotsByHost(db)

	req := httptest.NewRequest(http.MethodGet, "/hosts/125.199.235.74/snapshots", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(resp))
	}
}

// GET /diff
func TestGetDiff(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	id1 := insertSnapshot(t, db, "125.199.235.74", "2025-09-15T08:49:45Z", `{
		"timestamp":"2025-09-15T08:49:45Z","ip":"125.199.235.74",
		"services":[{"port":80,"protocol":"HTTP","status":301}]
	}`)

	id2 := insertSnapshot(t, db, "125.199.235.74", "2025-09-16T08:49:45Z", `{
		"timestamp":"2025-09-16T08:49:45Z","ip":"125.199.235.74",
		"services":[{"port":80,"protocol":"HTTP","status":200},
		            {"port":443,"protocol":"HTTPS","status":200}]
	}`)

	handler := GetDiff(db)

	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/diff?snap1=%d&snap2=%d", id1, id2), nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if len(resp["changed"].([]any)) != 1 {
		t.Errorf("expected 1 changed service, got %d", len(resp["changed"].([]any)))
	}
	if len(resp["added"].([]any)) != 1 {
		t.Errorf("expected 1 added service, got %d", len(resp["added"].([]any)))
	}
}
