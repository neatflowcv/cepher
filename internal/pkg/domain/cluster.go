package domain

import (
	"net"
	"strconv"
	"time"
)

type Cluster struct {
	id          string
	name        string
	hosts       []string
	key         string
	status      ClusterStatus
	lastBadTime time.Time
	detail      string
}

func NewCluster(
	id string,
	name string,
	hosts []string,
	key string,
	status ClusterStatus,
	lastBadTime time.Time,
	detail string,
) (*Cluster, error) {
	ret := Cluster{
		id:          id,
		name:        name,
		hosts:       hosts,
		key:         key,
		status:      status,
		lastBadTime: lastBadTime,
		detail:      detail,
	}

	err := ret.validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Cluster) SetStatus(status ClusterStatus, detail string, now time.Time) (*Cluster, error) {
	if now.After(c.lastBadTime) {
		return nil, InvalidParameterError("lastBadTime")
	}

	ret := c.clone()

	ret.status = status
	if !status.isHealthy() {
		ret.lastBadTime = now
		ret.detail = detail
	} else {
		ret.detail = ""
	}

	err := ret.validate()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Cluster) SetHosts(hosts []string) (*Cluster, error) {
	ret := c.clone()
	ret.hosts = hosts

	err := ret.validate()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Cluster) ID() string {
	return c.id
}

func (c *Cluster) Name() string {
	return c.name
}

func (c *Cluster) Hosts() []string {
	return c.hosts
}

func (c *Cluster) Key() string {
	return c.key
}

func (c *Cluster) Status() ClusterStatus {
	return c.status
}

func (c *Cluster) LastBadTime() time.Time {
	return c.lastBadTime
}

func (c *Cluster) Detail() string {
	return c.detail
}

func (c *Cluster) validate() error {
	if c.id == "" {
		return InvalidParameterError("id")
	}

	if c.name == "" {
		return InvalidParameterError("name")
	}

	if len(c.hosts) == 0 {
		return InvalidParameterError("hosts")
	}

	for _, host := range c.hosts {
		err := validateAddress(host)
		if err != nil {
			return err
		}
	}

	if c.key == "" {
		return InvalidParameterError("key")
	}

	err := c.status.validate()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) clone() *Cluster {
	return &Cluster{
		id:          c.id,
		name:        c.name,
		hosts:       c.hosts,
		key:         c.key,
		status:      c.status,
		lastBadTime: c.lastBadTime,
		detail:      c.detail,
	}
}

func validateAddress(address string) error {
	if address == "" {
		return InvalidParameterError("address")
	}

	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return InvalidParameterError("address")
	}

	if net.ParseIP(host) == nil {
		return InvalidParameterError("address")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return InvalidParameterError("address")
	}

	return nil
}
