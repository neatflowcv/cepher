package cephcli

type HealthDetail struct {
	Status string           `json:"status,omitempty"`
	Checks map[string]Check `json:"checks,omitempty"`
	Mutes  []any            `json:"mutes,omitempty"`
}

type Summary struct {
	Message string `json:"message,omitempty"`
	Count   int    `json:"count,omitempty"`
}

type Check struct {
	Severity string  `json:"severity,omitempty"`
	Summary  Summary `json:"summary"`
	Detail   []any   `json:"detail,omitempty"`
	Muted    bool    `json:"muted,omitempty"`
}
