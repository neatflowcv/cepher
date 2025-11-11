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
	Detail      any
}

func NewCluster(cluster *domain.Cluster) *Cluster {
	var hosts []string
	for _, host := range cluster.Hosts() {
		hosts = append(hosts, host.String())
	}

	return &Cluster{
		ID:          cluster.ID(),
		Name:        cluster.Name(),
		Hosts:       hosts,
		Key:         cluster.Key(),
		Status:      string(cluster.Status()),
		LastBadTime: cluster.LastBadTime(),
		Detail:      cluster.Detail(),
	}
}

func (c *Cluster) ToDomain() (*domain.Cluster, error) {
	addresses, err := domain.NewAddressesFromHosts(c.Hosts)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain addresses: %w", err)
	}

	cluster, err := domain.NewCluster(
		c.ID, c.Name, addresses, c.Key, domain.ClusterStatus(c.Status), c.LastBadTime, c.Detail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain cluster: %w", err)
	}

	return cluster, nil
}
