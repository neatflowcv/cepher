package flow

import "time"

type Cluster struct {
	ID       string
	Name     string
	Status   string
	IsStable bool
}

type RegisterCluster struct {
	Name  string
	Hosts []string
	Key   string
	Now   time.Time
}
