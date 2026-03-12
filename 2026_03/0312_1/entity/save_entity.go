package entity

import (
	"time"
)

type SaveEntity struct {
	Timer *time.Timer
	Msg   []*Message
}
