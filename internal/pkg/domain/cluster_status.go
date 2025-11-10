package domain

type ClusterStatus string

const (
	ClusterStatusUnknown       ClusterStatus = "HEALTH_UNKNOWN"
	ClusterStatusHealthOK      ClusterStatus = "HEALTH_OK"
	ClusterStatusHealthWarning ClusterStatus = "HEALTH_WARN"
	ClusterStatusHealthError   ClusterStatus = "HEALTH_ERR"
)

func (s ClusterStatus) isHealthy() bool {
	return s == ClusterStatusHealthOK
}

func (s ClusterStatus) validate() error {
	switch s {
	case ClusterStatusUnknown,
		ClusterStatusHealthOK,
		ClusterStatusHealthWarning,
		ClusterStatusHealthError:
		return nil
	default:
		return InvalidParameterError("status")
	}
}
