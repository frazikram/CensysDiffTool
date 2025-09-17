package models

type Snapshot struct {
	Timestamp    string    `json:"timestamp"`
	IP           string    `json:"ip"`
	Services     []Service `json:"services"`
	ServiceCount int       `json:"service_count"`
}

type Service struct {
	Port            int       `json:"port"`
	Protocol        string    `json:"protocol"`
	Status          int       `json:"status"`
	Software        *Software `json:"software,omitempty"`
	TLS             *TLSInfo  `json:"tls,omitempty"`
	Vulnerabilities []string  `json:"vulnerabilities,omitempty"`
}
type Software struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
	Version string `json:"version,omitempty"`
}

type TLSInfo struct {
	Version               string `json:"version"`
	Cipher                string `json:"cipher"`
	CertFingerprintSHA256 string `json:"cert_fingerprint_sha256"`
}
