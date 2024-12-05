package stats

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

// prometheusClient implements Client interface to expose metric for prometheus server.
type prometheusClient struct {
	registerer prometheus.Registerer
	counter    *counter
	histogram  *histogram
	gauge      *gauge
}

func NewPrometheusClient() Client {
	return NewPrometheusClientWithRegisterer(prometheus.DefaultRegisterer)

}

func NewPrometheusClientWithRegisterer(registerer prometheus.Registerer) Client {
	return &prometheusClient{
		registerer: registerer,
		counter:    newCounter(registerer),
		histogram:  newHistogram(registerer),
		gauge:      newGauge(registerer),
	}
}

func (c *prometheusClient) CountInc(ctx context.Context, contextTag, key string, tags ...Tag) {
	if err := c.counter.Add(contextTag, key, 1, tags...); err != nil {
		logrus.Warnf("cannot do increase with contextTag: %s, key: %s, tags: %v, cause: %v", contextTag, key, tags, err)
	}
}

func (c *prometheusClient) Count(ctx context.Context, contextTag, key string, amount int64, tags ...Tag) {
	if err := c.counter.Add(contextTag, key, amount, tags...); err != nil {
		logrus.Warnf("cannot update count with contextTag: %s, key: %s, tags: %v, amount: %d, cause: %v", contextTag, key, tags, amount, err)
	}
}

func (c *prometheusClient) Histogram(ctx context.Context, contextTag, key string, value float64, tags ...Tag) {
	if err := c.histogram.Observe(contextTag, key, value, tags...); err != nil {
		logrus.Warnf("cannot update histogram with contextTag: %s, key: %s, tags: %v, value: %f, cause: %v", contextTag, key, tags, value, err)
	}
}

func (c *prometheusClient) Duration(ctx context.Context, contextTag, key string, start time.Time, tags ...Tag) {
	elapsed := float64(time.Since(start)) / float64(time.Second)
	if err := c.histogram.Observe(contextTag, key, elapsed, tags...); err != nil {
		logrus.Warnf("cannot update duration with contextTag: %s, key: %s, tags: %v, duration: %f, cause: %v", contextTag, key, tags, elapsed, err)
	}
}

func (c *prometheusClient) Gauge(ctx context.Context, contextTag, key string, value float64, tags ...Tag) {
	if err := c.gauge.Set(contextTag, key, value, tags...); err != nil {
		logrus.Warnf("cannot update gauge with contextTag: %s, key: %s, tags: %v, value: %f, cause: %v", contextTag, key, tags, value, err)
	}
}

func (c *prometheusClient) GaugeInc(ctx context.Context, contextTag, key string, tags ...Tag) {
	if err := c.gauge.Inc(contextTag, key, tags...); err != nil {
		logrus.Warnf("cannot increase gauge with contextTag: %s, key: %s, tags: %v, cause: %v", contextTag, key, tags, err)
	}
}

func (c *prometheusClient) GaugeDec(ctx context.Context, contextTag, key string, tags ...Tag) {
	if err := c.gauge.Dec(contextTag, key, tags...); err != nil {
		logrus.Warnf("cannot decrease gauge with contextTag: %s, key: %s, tags: %v, cause: %v", contextTag, key, tags, err)
	}
}
