package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/neatflowcv/cepher/internal/pkg/domain"
	"github.com/neatflowcv/cepher/internal/pkg/repository"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	path string
}

func NewRepository(path string) *Repository {
	return &Repository{
		path: path,
	}
}

func (r *Repository) CreateCluster(ctx context.Context, dCluster *domain.Cluster) error {
	cluster := NewCluster(dCluster)
	path := filepath.Join(r.path, dCluster.ID())
	filePath := path + ".json"

	_, err := os.Stat(filePath)
	if err == nil {
		return repository.ErrClusterAlreadyExists
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	data, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cluster: %w", err)
	}

	const permission = 0600

	err = os.WriteFile(filePath, data, permission)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *Repository) ListClusters(ctx context.Context) ([]*domain.Cluster, error) {
	files, err := os.ReadDir(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var clusters []*Cluster

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if len(name) < 6 || name[len(name)-5:] != ".json" {
			continue
		}

		filePath := filepath.Clean(filepath.Join(r.path, name))

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		var cluster Cluster

		err = json.Unmarshal(data, &cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal cluster file %s: %w", filePath, err)
		}

		clusters = append(clusters, &cluster)
	}

	var ret []*domain.Cluster

	for _, cluster := range clusters {
		dCluster, err := cluster.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert cluster to domain: %w", err)
		}

		ret = append(ret, dCluster)
	}

	return ret, nil
}
