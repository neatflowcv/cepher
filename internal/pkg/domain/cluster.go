package domain

import (
	"reflect"
	"time"
)

type Cluster struct {
	id          string
	name        string
	hosts       []*Address
	key         string
	status      ClusterStatus
	lastBadTime time.Time
	detail      any
}

func NewCluster(
	id string,
	name string,
	hosts []*Address,
	key string,
	status ClusterStatus,
	lastBadTime time.Time,
	detail any,
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

func (c *Cluster) SetStatus(status ClusterStatus, detail any, now time.Time) (*Cluster, error) {
	if now.Before(c.lastBadTime) {
		return nil, InvalidParameterError("lastBadTime")
	}

	lastBadTime := c.lastBadTime
	if !status.isHealthy() {
		lastBadTime = now
	}

	if c.status == status &&
		c.lastBadTime.Equal(lastBadTime) &&
		reflect.DeepEqual(c.detail, detail) {
		return c, nil
	}

	ret := c.clone()

	ret.status = status
	ret.lastBadTime = lastBadTime
	ret.detail = detail

	err := ret.validate()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Cluster) SetHosts(hosts []*Address) (*Cluster, error) {
	if reflect.DeepEqual(c.hosts, hosts) {
		return c, nil
	}

	ret := c.clone()
	ret.hosts = hosts

	return ret, nil
}

func (c *Cluster) IsOK() bool {
	return c.status.isHealthy()
}

func (c *Cluster) ID() string {
	return c.id
}

func (c *Cluster) Name() string {
	return c.name
}

func (c *Cluster) Hosts() []*Address {
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

func (c *Cluster) Detail() any {
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
		if host == nil {
			return InvalidParameterError("hosts")
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
