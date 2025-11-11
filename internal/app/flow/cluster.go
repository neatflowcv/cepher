package flow

import (
	"time"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
)

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

func NewCluster(cluster *domain.Cluster) *Cluster {
	return &Cluster{
		ID:       cluster.ID(),
		Name:     cluster.Name(),
		Status:   string(cluster.Status()),
		IsStable: domain.IsClusterStable(cluster, cluster.LastBadTime()),
	}
}

func NewClusters(clusters []*domain.Cluster) []*Cluster {
	var ret []*Cluster
	for _, cluster := range clusters {
		ret = append(ret, NewCluster(cluster))
	}

	return ret
}
