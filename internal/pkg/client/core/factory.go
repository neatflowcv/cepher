package core

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/neatflowcv/cepher/internal/pkg/client"
	"github.com/neatflowcv/cepher/internal/pkg/domain"
	"github.com/neatflowcv/cepher/pkg/cephsetup"
)

var _ client.Factory = (*Factory)(nil)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewClient(ctx context.Context, cluster *domain.Cluster) (client.Client, error) {
	tempDir, err := os.MkdirTemp("", "cepher")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	hosts := cluster.Hosts()
	rand.Shuffle(len(hosts), func(i, j int) {
		hosts[i], hosts[j] = hosts[j], hosts[i]
	})

	var hostStrings []string
	for _, host := range hosts {
		hostStrings = append(hostStrings, host.String())
	}

	err = cephsetup.Setup(tempDir, hostStrings, cluster.Key())
	if err != nil {
		return nil, fmt.Errorf("failed to setup ceph: %w", err)
	}

	return newClient(tempDir, "20.1.1"), nil
}
