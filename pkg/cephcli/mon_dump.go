package cephcli

import "time"

type MonDump struct {
	Epoch             int       `json:"epoch,omitempty"`
	Fsid              string    `json:"fsid,omitempty"`
	Modified          time.Time `json:"modified,omitzero"`
	Created           time.Time `json:"created,omitzero"`
	MinMonRelease     int       `json:"min_mon_release,omitempty"`
	MinMonReleaseName string    `json:"min_mon_release_name,omitempty"`
	ElectionStrategy  int       `json:"election_strategy,omitempty"`
	DisallowedLeaders string    `json:"disallowed_leaders,omitempty"`
	StretchMode       bool      `json:"stretch_mode,omitempty"`
	TiebreakerMon     string    `json:"tiebreaker_mon,omitempty"`
	RemovedRanks      string    `json:"removed_ranks,omitempty"`
	Features          Features  `json:"features"`
	Mons              []Mons    `json:"mons,omitempty"`
	Quorum            []int     `json:"quorum,omitempty"`
}

type Features struct {
	Persistent []string `json:"persistent,omitempty"`
	Optional   []any    `json:"optional,omitempty"`
}

type Addrvec struct {
	Type  string `json:"type,omitempty"`
	Addr  string `json:"addr,omitempty"`
	Nonce int    `json:"nonce,omitempty"`
}

type PublicAddrs struct {
	Addrvec []Addrvec `json:"addrvec,omitempty"`
}

type Mons struct {
	Rank          int         `json:"rank,omitempty"`
	Name          string      `json:"name,omitempty"`
	PublicAddrs   PublicAddrs `json:"public_addrs"`
	Addr          string      `json:"addr,omitempty"`
	PublicAddr    string      `json:"public_addr,omitempty"`
	Priority      int         `json:"priority,omitempty"`
	Weight        int         `json:"weight,omitempty"`
	CrushLocation string      `json:"crush_location,omitempty"`
}
