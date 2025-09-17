package services

import (
	"fmt"
	"hostdiff/api/internal/models"
)

//Compares two snapshots and returns diff

func DiffSnapshots(snap1, snap2 models.Snapshot) models.DiffResult {
	result := models.DiffResult{}

	//Building a map for easier Comparision
	map1 := make(map[string]models.Service)
	for _, svc := range snap1.Services {
		key := serviceKey(svc)
		map1[key] = svc
	}

	map2 := make(map[string]models.Service)
	for _, svc := range snap2.Services {
		key := serviceKey(svc)
		map2[key] = svc
	}
	//See what was removed and changed
	for k, svc1 := range map1 {
		if svc2, exists := map2[k]; !exists {
			// In snap1 but not snap2 → removed
			result.Removed = append(result.Removed, svc1)
		} else {
			// In both → check if changed
			if serviceChanged(svc1, svc2) {
				result.Changed = append(result.Changed, models.ChangeSet{
					Port:     svc1.Port,
					Protocol: svc1.Protocol,
					From:     svc1,
					To:       svc2,
				})
			}
		}
	}
	// See what was added
	for k, svc2 := range map2 {
		if _, exists := map1[k]; !exists {
			result.Added = append(result.Added, svc2)
		}
	}
	return result
}

// Creating a service key for uniqueness
func serviceKey(s models.Service) string {
	return fmt.Sprintf("%d/%s", s.Port, s.Protocol)
}

// serviceChanged checks if any fields differ
func serviceChanged(a, b models.Service) bool {
	if a.Status != b.Status {
		return true
	}
	if (a.Software != nil && b.Software != nil) && *a.Software != *b.Software {
		return true
	}
	if (a.TLS != nil && b.TLS != nil) && *a.TLS != *b.TLS {
		return true
	}
	if len(a.Vulnerabilities) != len(b.Vulnerabilities) {
		return true
	}
	// TODO: deeper vuln comparison if needed in the future
	return false
}
