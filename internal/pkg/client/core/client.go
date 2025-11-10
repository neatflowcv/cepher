package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/neatflowcv/cepher/internal/pkg/client"
	"github.com/neatflowcv/cepher/internal/pkg/domain"
	"github.com/neatflowcv/cepher/pkg/cephcli"
)

var _ client.Client = (*Client)(nil)

type Client struct {
	client *cephcli.Client
	path   string
}

func newClient(path, version string) *Client {
	return &Client{
		client: cephcli.NewClient(path, version),
		path:   path,
	}
}

func (c *Client) Close() {
	err := os.RemoveAll(c.path)
	if err != nil {
		log.Printf("failed to remove temporary directory: %v", err)
	}
}

func (c *Client) HealthCheck(ctx context.Context) (domain.ClusterStatus, error) {
	health, err := c.client.GetHealth(ctx)
	if err != nil {
		return domain.ClusterStatusUnknown, fmt.Errorf("failed to get health: %w", err)
	}

	status := strings.Split(health, " ")[0]
	status = strings.TrimSpace(status)
	log.Println("status", status)
	log.Println("health", health)

	return domain.ClusterStatus(status), nil
}
