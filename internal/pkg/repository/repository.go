package repository

import (
	"context"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
)

type Repository interface {
	CreateCluster(ctx context.Context, cluster *domain.Cluster) error
}
