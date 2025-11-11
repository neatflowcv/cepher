package client

import (
	"context"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
)

type Factory interface {
	NewClient(ctx context.Context, cluster *domain.Cluster) (Client, error)
}

type Client interface {
	Close()

	HealthCheck(ctx context.Context) (domain.ClusterStatus, any, error)
	ListMonitors(ctx context.Context) ([]*domain.Address, error)
}
