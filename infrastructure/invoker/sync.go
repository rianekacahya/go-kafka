package invoker

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"time"
)

type sync struct {
	redisgo *redis.Client
}

func SyncInitialize(redisgo *redis.Client) *sync {
	return &sync{redisgo}
}

func (s *sync) EventWhisper(ctx context.Context, key string) error {
	if _, err := s.redisgo.RPush(key, nil).Result(); err != nil {
		return err
	}

	return nil
}

func (s *sync) EventListener(ctx context.Context, key string) <-chan error {
	errx := make(chan error)
	go func() {
		defer close(errx)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, err := s.redisgo.BRPop(0*time.Second, key).Result(); err != nil {
					errx <- crashy.Wrap(err, crashy.ErrCodeNetConnect, "error when listening to whisper service")
					return
				}

				return
			}
		}
	}()

	return errx
}
