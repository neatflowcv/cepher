package flow

import (
	"context"
	"fmt"

	"github.com/neatflowcv/cepher/internal/pkg/client"
	"github.com/neatflowcv/cepher/internal/pkg/domain"
	"github.com/neatflowcv/cepher/internal/pkg/idgenerator"
	"github.com/neatflowcv/cepher/internal/pkg/repository"
)

type Service struct {
	idGenerator idgenerator.Generator
	factory     client.Factory
	repository  repository.Repository
}

func NewService(idGenerator idgenerator.Generator, factory client.Factory, repository repository.Repository) *Service {
	return &Service{
		idGenerator: idGenerator,
		factory:     factory,
		repository:  repository,
	}
}

func (s *Service) RegisterCluster(ctx context.Context, registerCluster *RegisterCluster) (*Cluster, error) {
	id := s.idGenerator.GenerateID()

	cluster, err := domain.NewCluster(
		id, registerCluster.Name, registerCluster.Hosts, registerCluster.Key,
		domain.ClusterStatusUnknown, registerCluster.Now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	client, err := s.factory.NewClient(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	status, err := client.HealthCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to health check: %w", err)
	}

	cluster, err = cluster.SetStatus(status, registerCluster.Now)
	if err != nil {
		return nil, fmt.Errorf("failed to set status: %w", err)
	}

	err = s.repository.CreateCluster(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	return &Cluster{
		ID:       cluster.ID(),
		Name:     cluster.Name(),
		Status:   string(cluster.Status()),
		IsStable: domain.IsClusterStable(cluster, registerCluster.Now),
	}, nil
}
