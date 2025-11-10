package cephcli

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	path    string
	version string
}

func NewClient(path, version string) *Client {
	return &Client{
		path:    path,
		version: version,
	}
}

func (c *Client) GetHealth(ctx context.Context) (string, error) {
	image := "quay.io/ceph/ceph:v" + c.version
	volume := c.path + ":/etc/ceph"
	cmd := exec.CommandContext( //nolint:gosec
		ctx,
		"podman", "run", "--rm", "-v", volume, image, "ceph", "health", "detail",
	)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
