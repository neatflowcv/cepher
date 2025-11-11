package core

import (
	"context"
	"fmt"
	"log"
	"os"

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

func (c *Client) HealthCheck(ctx context.Context) (domain.ClusterStatus, any, error) {
	health, err := c.client.HealthDetail(ctx)
	if err != nil {
		return domain.ClusterStatusUnknown, "", fmt.Errorf("failed to get health: %w", err)
	}

	return domain.ClusterStatus(health.Status), health.Checks, nil
}

func (c *Client) ListMonitors(ctx context.Context) ([]*domain.Address, error) {
	dump, err := c.client.MonDump(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitors: %w", err)
	}

	var ret []*domain.Address

	for _, mon := range dump.Mons {
		var (
			maxAddr string
			maxType string
		)

		for _, addr := range mon.PublicAddrs.Addrvec {
			if addr.Type > maxType {
				maxType = addr.Type
				maxAddr = addr.Addr
			}
		}

		address, err := domain.NewAddressFromHost(maxAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain address: %w", err)
		}

		ret = append(ret, address)
	}

	return ret, nil
}
