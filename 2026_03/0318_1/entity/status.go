package entity

import "time"

type Status struct {
	Service   string
	Latency   time.Duration
	Timestamp time.Time
}
