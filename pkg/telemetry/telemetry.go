package telemetry

import (
	"context"
	newrelic "github.com/newrelic/go-agent"
)

const (
	ContextKeyTelemetry = "telemetry"
)

func New(service, key string) (newrelic.Application, error) {
	config := newrelic.NewConfig(service, key)
	config.DistributedTracer.Enabled = true
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func ContextTelemetry(ctx context.Context) newrelic.Transaction {
	if nr, ok := ctx.Value(ContextKeyTelemetry).(newrelic.Transaction); ok {
		return nr
	}
	return nil
}
