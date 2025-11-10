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
