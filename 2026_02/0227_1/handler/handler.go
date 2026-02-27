package handler

import (
	"chat/entity"
	"context"
)

// handler라는 함수 타입으로 사용함.
type Handler func(ctx context.Context, msg entity.Message, idx int) error
