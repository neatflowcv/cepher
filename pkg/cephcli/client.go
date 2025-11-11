package cephcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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

func (c *Client) HealthDetail(ctx context.Context) (*HealthDetail, error) {
	image := "quay.io/ceph/ceph:v" + c.version
	volume := c.path + ":/etc/ceph"
	cmd := exec.CommandContext( //nolint:gosec
		ctx,
		"podman", "run", "--rm", "-v", volume, image, "ceph", "health", "detail", "-f", "json",
	)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w: %s", err, stderr.String())
	}

	var ret HealthDetail

	err = json.NewDecoder(&stdout).Decode(&ret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode health detail: %w", err)
	}

	return &ret, nil
}

func (c *Client) MonDump(ctx context.Context) (*MonDump, error) {
	image := "quay.io/ceph/ceph:v" + c.version
	volume := c.path + ":/etc/ceph"
	cmd := exec.CommandContext( //nolint:gosec
		ctx,
		"podman", "run", "--rm", "-v", volume, image, "ceph", "mon", "dump", "-f", "json",
	)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w: %s", err, stderr.String())
	}

	var ret MonDump

	err = json.NewDecoder(&stdout).Decode(&ret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mon dump: %w", err)
	}

	return &ret, nil
}
