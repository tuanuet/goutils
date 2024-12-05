// Package stats provides methods to communicate with prometheus via stats interface.
package stats

import (
	"context"
	"time"
)

// Client defines the statsD prometheusClient exported interface.
type Client interface {
	// CountInc increase the metric by 1.
	CountInc(ctx context.Context, contextTag, key string, tags ...Tag)
	// Count increase the metric by amount.
	Count(ctx context.Context, contextTag, key string, amount int64, tags ...Tag)
	// Histogram set value for the histogram metric.
	Histogram(ctx context.Context, contextTag, key string, value float64, tags ...Tag)
	// Duration set duration from start moment to the metric
	Duration(ctx context.Context, contextTag, key string, start time.Time, tags ...Tag)
	// Gauge set value for gauge
	Gauge(ctx context.Context, contextTag, key string, value float64, tags ...Tag)
	// GaugeInc increments value of gauge by 1.
	GaugeInc(ctx context.Context, contextTag, key string, tags ...Tag)
	// GaugeDec decrements value of gauge by 1.
	GaugeDec(ctx context.Context, contextTag, key string, tags ...Tag)
}

type Tag struct {
	Key   string
	Value string
}

//go:generate mockery --case underscore --name=Client --inpackage --filename=client_mock.go
