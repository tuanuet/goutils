package stats

import (
	"context"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func getMetric(col prometheus.Collector) dto.Metric {
	c := make(chan prometheus.Metric, 1)
	col.Collect(c)
	m := dto.Metric{}
	_ = (<-c).Write(&m)
	return m
}

func TestClient(t *testing.T) {
	tests := []struct {
		name            string
		contextTag, key string
		tags            []Tag
		operation       func(c *prometheusClient, contextTag, key string, tags ...Tag)
		getValue        func(c *prometheusClient, contextTag, key string, tags ...Tag) float64
		expectedValue   float64
		panic           string
	}{
		{
			name:       "Test CountInc func without contextTag",
			contextTag: "",
			key:        "count1",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.CountInc(context.Background(), contextTag, key, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.counter.GetOrCreateCollector(contextTag, key, tags...))
				return *metric.Counter.Value
			},
			expectedValue: float64(1),
		},
		{
			name:       "Test Count func without Tags",
			contextTag: "svc",
			key:        "count1",
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.Count(context.Background(), contextTag, key, 12, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.counter.GetOrCreateCollector(contextTag, key, tags...))
				return *metric.Counter.Value
			},
			expectedValue: float64(12),
		},
		{
			name:       "Test Histogram func in sum value",
			contextTag: "",
			key:        "histogram1",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.Histogram(context.Background(), contextTag, key, 7, tags...)
				c.Histogram(context.Background(), contextTag, key, -9, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.histogram.GetOrCreateCollector(contextTag, key, tags...))
				return *metric.Histogram.SampleSum
			},
			expectedValue: float64(-2),
		},
		{
			name:       "Test Histogram func in count value",
			contextTag: "",
			key:        "histogram1",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.Histogram(context.Background(), contextTag, key, 7, tags...)
				c.Histogram(context.Background(), contextTag, key, -9, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.histogram.GetOrCreateCollector(contextTag, key, tags...))
				return float64(*metric.Histogram.SampleCount)
			},
			expectedValue: float64(2),
		},
		{
			name:       "Test registry two goutils have different type but same name",
			contextTag: "svc",
			key:        "metric_name",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.Histogram(context.Background(), contextTag, key, 1.3, tags...)
				c.CountInc(context.Background(), contextTag, key, tags...)
			},
			panic: "duplicate goutils collector registration attempted",
		},
		{
			name:       "Test Count an negative number",
			contextTag: "svc",
			key:        "metric_name",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.Count(context.Background(), contextTag, key, -2, tags...)
			},
			panic: "counter cannot decrease in value",
		},
		{
			name:       "Test Gauge Increase",
			contextTag: "svc",
			key:        "metric_name",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.GaugeInc(context.Background(), contextTag, key, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.gauge.GetOrCreateCollector(contextTag, key, tags...))
				return *metric.Gauge.Value
			},
			expectedValue: float64(1),
		},
		{
			name:       "Test Gauge Decrease",
			contextTag: "svc",
			key:        "metric_name",
			tags:       []Tag{{Key: "key1", Value: "value1"}},
			operation: func(c *prometheusClient, contextTag, key string, tags ...Tag) {
				c.GaugeDec(context.Background(), contextTag, key, tags...)
			},
			getValue: func(c *prometheusClient, contextTag, key string, tags ...Tag) float64 {
				metric := getMetric(c.gauge.GetOrCreateCollector(contextTag, key, tags...))
				return *metric.Gauge.Value
			},
			expectedValue: float64(-1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := prometheus.NewRegistry()
			client := NewPrometheusClientWithRegisterer(registry).(*prometheusClient)
			if len(tt.panic) > 0 {
				assert.PanicsWithValue(t, tt.panic, func() {
					defer func() {
						r := recover()
						if r != nil {
							panic(fmt.Sprintf("%v", r))
						}
					}()
					tt.operation(client, tt.contextTag, tt.key, tt.tags...)
				})
			} else {
				tt.operation(client, tt.contextTag, tt.key, tt.tags...)
				actual := tt.getValue(client, tt.contextTag, tt.key, tt.tags...)
				assert.Equal(t, tt.expectedValue, actual)
			}
		})
	}
}
