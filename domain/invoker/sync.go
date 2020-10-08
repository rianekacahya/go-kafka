package invoker

import (
	"context"
)

type SyncInvoker interface {
	EventWhisper(context.Context, string) error
	EventListener(context.Context, string) <-chan error
}
