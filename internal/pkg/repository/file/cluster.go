package file

import (
	"time"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
)

type Cluster struct {
	ID          string
	Name        string
	Hosts       []string
	Key         string
	Status      string
	LastBadTime time.Time
	Detail      string
}

func NewCluster(cluster *domain.Cluster) *Cluster {
	return &Cluster{
		ID:          cluster.ID(),
		Name:        cluster.Name(),
		Hosts:       cluster.Hosts(),
		Key:         cluster.Key(),
		Status:      string(cluster.Status()),
		LastBadTime: cluster.LastBadTime(),
		Detail:      cluster.Detail(),
	}
}
