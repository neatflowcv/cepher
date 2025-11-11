package repository

import (
	"context"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
)

type Repository interface {
	CreateCluster(ctx context.Context, cluster *domain.Cluster) error
	ListClusters(ctx context.Context) ([]*domain.Cluster, error)
	GetCluster(ctx context.Context, id string) (*domain.Cluster, error)
	UpdateCluster(ctx context.Context, cluster *domain.Cluster) error
}
