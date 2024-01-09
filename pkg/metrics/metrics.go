package metrics

import "context"

type Metrics interface {
	RegisterCounter(name string) error
	RegisterGauge(name string) error
	Set(name string, value int64)
	Add(name string, value int64)
	Inc(name string)
	Get(name string) int64
}

// NewMetrics always returns Screener for now.
func NewMetrics(ctx context.Context) Metrics {
	return NewScreener(ctx)
}
