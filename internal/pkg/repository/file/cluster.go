package file

import (
	"fmt"
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

func (c *Cluster) ToDomain() (*domain.Cluster, error) {
	cluster, err := domain.NewCluster(
		c.ID, c.Name, c.Hosts, c.Key, domain.ClusterStatus(c.Status), c.LastBadTime, c.Detail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain cluster: %w", err)
	}

	return cluster, nil
}
