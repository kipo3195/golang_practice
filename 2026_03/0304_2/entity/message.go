package entity

import "time"

type Message struct {
	Value string
	Exp   time.Time
}
