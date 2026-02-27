package processor

import (
	"chat/entity"
	"context"
)

type Processor interface {
	Start(ctx context.Context)
	Submit(msg entity.Message) error
	Stop()
}
