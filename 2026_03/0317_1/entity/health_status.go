package entity

import "time"

type HealthStatus struct {
	ServiceName string
	IsAlive     bool
	Latency     time.Duration
}
