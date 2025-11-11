package flow

import (
	"context"
	"fmt"
	"log"
	"time"

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

	addresses, err := domain.NewAddressesFromHosts(registerCluster.Hosts)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain addresses: %w", err)
	}

	cluster, err := domain.NewCluster(
		id, registerCluster.Name, addresses, registerCluster.Key,
		domain.ClusterStatusUnknown, registerCluster.Now,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	client, err := s.factory.NewClient(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	status, detail, err := client.HealthCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to health check: %w", err)
	}

	cluster, err = cluster.SetStatus(status, detail, registerCluster.Now)
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
		Detail:   cluster.Detail(),
	}, nil
}

func (s *Service) ListClusters(ctx context.Context) ([]*Cluster, error) {
	clusters, err := s.repository.ListClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	return NewClusters(clusters), nil
}

// RefreshCluster refreshes the cluster status
// returns true if the cluster status is ok.
func (s *Service) RefreshCluster(ctx context.Context, id string, now time.Time) (bool, error) {
	cluster, err := s.repository.GetCluster(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to get cluster: %w", err)
	}

	client, err := s.factory.NewClient(ctx, cluster)
	if err != nil {
		return false, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	var changedCluster *domain.Cluster

	status, detail, err := client.HealthCheck(ctx)
	if err != nil {
		// client로 부터 상태를 가져오지 못하면, 상태를 Unknown으로 설정하고 계속 진행한다.
		log.Printf("failed to health check cluster %s: %v", id, err)

		status = domain.ClusterStatusUnknown
		detail = ""
	}

	changedCluster, err = cluster.SetStatus(status, detail, now)
	if err != nil {
		return false, fmt.Errorf("failed to set status: %w", err)
	}

	if cluster == changedCluster {
		return changedCluster.IsOK(), nil
	}

	err = s.repository.UpdateCluster(ctx, changedCluster)
	if err != nil {
		return false, fmt.Errorf("failed to update cluster: %w", err)
	}

	return changedCluster.IsOK(), nil
}
