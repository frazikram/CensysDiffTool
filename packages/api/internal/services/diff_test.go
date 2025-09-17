package services

import (
	"hostdiff/api/internal/models"
	"testing"
)

func TestDiffSnapshots(t *testing.T) {
	snap1 := models.Snapshot{
		IP: "125.199.235.74",
		Services: []models.Service{
			{Port: 80, Protocol: "HTTP", Status: 301,
				Software: &models.Software{Vendor: "microsoft", Product: "iis", Version: "8.5"}},
		},
	}

	snap2 := models.Snapshot{
		IP: "125.199.235.74",
		Services: []models.Service{
			{Port: 80, Protocol: "HTTP", Status: 200,
				Software: &models.Software{Vendor: "microsoft", Product: "iis", Version: "10.0"}},
			{Port: 443, Protocol: "HTTPS", Status: 200,
				Software: &models.Software{Vendor: "microsoft", Product: "asp.net"}},
		},
	}

	diff := DiffSnapshots(snap1, snap2)

	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed service, got %d", len(diff.Changed))
	}
	if len(diff.Added) != 1 {
		t.Errorf("expected 1 added service, got %d", len(diff.Added))
	}
	if len(diff.Removed) != 0 {
		t.Errorf("expected 0 removed services, got %d", len(diff.Removed))
	}
}
