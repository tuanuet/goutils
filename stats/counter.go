package stats

import (
	"github.com/prometheus/client_golang/prometheus"
)

type counter struct {
	*collector
}

func newCounter(registerer prometheus.Registerer) *counter {

	createCounterFn := func(contextTag, key string, tagKeys ...string) prometheus.Collector {
		return prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: contextTag,
				Name:      key,
			},
			tagKeys,
		)
	}

	return &counter{
		collector: newCollector(registerer, createCounterFn),
	}
}

// Add to counter with amount value.
func (c *counter) Add(contextTag, key string, amount int64, tags ...Tag) error {
	counter, err := c.getOrCreateCounter(contextTag, key, tags...)
	if err != nil {
		return err
	}
	counter.Add(float64(amount))
	return nil
}

func (c *counter) getOrCreateCounter(contextTag, key string, tags ...Tag) (prometheus.Counter, error) {
	collector := c.collector.GetOrCreateCollector(contextTag, key, tags...)
	counterVec := collector.(*prometheus.CounterVec)

	tagMap := getTagMap(tags...)
	return counterVec.GetMetricWith(tagMap)
}
