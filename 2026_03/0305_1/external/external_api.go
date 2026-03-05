package external

import (
	"context"
	"math/rand"
	"test/entity"
	"time"
)

// 가상의 외부 알림 API
func SendPushNotification(ctx context.Context, msg entity.Message) error {
	select {
	case <-time.After(time.Duration(rand.Intn(800)) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
