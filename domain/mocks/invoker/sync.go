package invoker

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Sync struct {
	mock.Mock
}

func (s *Sync) EventWhisper(ctx context.Context, key string) error {
	ret := s.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (s *Sync) EventListener(ctx context.Context, key string) <- chan error {
	var errx = make(chan error)
	close(errx)
	return errx
}
