package stats

import (
	"context"
	"time"
)

// StubClient just a holder
type StubClient struct {
}

// Count provides a mock function with given fields: ctx, contextTag, key, amount, tags
func (_m *StubClient) Count(ctx context.Context, contextTag string, key string, amount int64, tags ...Tag) {
}

// CountInc provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *StubClient) CountInc(ctx context.Context, contextTag string, key string, tags ...Tag) {
}

// Duration provides a mock function with given fields: ctx, contextTag, key, start, tags
func (_m *StubClient) Duration(ctx context.Context, contextTag string, key string, start time.Time, tags ...Tag) {
}

// Gauge provides a mock function with given fields: ctx, contextTag, key, value, tags
func (_m *StubClient) Gauge(ctx context.Context, contextTag string, key string, value float64, tags ...Tag) {
}

// GaugeDec provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *StubClient) GaugeDec(ctx context.Context, contextTag string, key string, tags ...Tag) {
}

// GaugeInc provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *StubClient) GaugeInc(ctx context.Context, contextTag string, key string, tags ...Tag) {
}

// Histogram provides a mock function with given fields: ctx, contextTag, key, value, tags
func (_m *StubClient) Histogram(ctx context.Context, contextTag string, key string, value float64, tags ...Tag) {
}
