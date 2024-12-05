package stats

import (
	"github.com/prometheus/client_golang/prometheus"
)

type histogram struct {
	*collector
}

func newHistogram(registerer prometheus.Registerer) *histogram {
	createHistogramFn := func(contextTag, key string, tagKeys ...string) prometheus.Collector {
		return prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: contextTag,
				Name:      key,
			}, tagKeys,
		)
	}
	return &histogram{
		collector: newCollector(registerer, createHistogramFn),
	}
}

// Observe histogram with the value.
func (h *histogram) Observe(contextTag, key string, value float64, tags ...Tag) error {
	histogram, err := h.getOrCreateHistogram(contextTag, key, tags...)
	if err != nil {
		return err
	}
	histogram.Observe(value)
	return nil
}

func (h *histogram) getOrCreateHistogram(contextTag, key string, tags ...Tag) (prometheus.Observer, error) {
	collector := h.collector.GetOrCreateCollector(contextTag, key, tags...)
	histogramVec := collector.(*prometheus.HistogramVec)

	tagMap := getTagMap(tags...)
	return histogramVec.GetMetricWith(tagMap)
}
