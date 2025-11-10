package domain

import "time"

func IsClusterStable(cluster *Cluster, now time.Time) bool {
	const stableInterval = 3 * time.Minute

	return cluster.LastBadTime().Add(stableInterval).Before(now)
}
