package models

type DiffResult struct {
	Added   []Service   `json:"added"`
	Removed []Service   `json:"removed"`
	Changed []ChangeSet `json:"changed"`
}

type ChangeSet struct {
	Port     int     `json:"port"`
	Protocol string  `json:"protocol"`
	From     Service `json:"from"`
	To       Service `json:"to"`
}
