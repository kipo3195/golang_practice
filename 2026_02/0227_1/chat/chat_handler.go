package chat

import (
	"chat/entity"
	"chat/handler"
	"context"
	"log"
)

func NewChatHandler() handler.Handler {

	return func(ctx context.Context, msg entity.Message, idx int) error {
		log.Printf("실행하는 고루틴 idx : %d, msg: %s", idx, msg)
		return nil
	}
}
