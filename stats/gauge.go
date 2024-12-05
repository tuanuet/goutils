package stats

import (
	"github.com/prometheus/client_golang/prometheus"
)

type gauge struct {
	*collector
}

func newGauge(registerer prometheus.Registerer) *gauge {
	createGaugeFn := func(contextTag, key string, tagKeys ...string) prometheus.Collector {
		return prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: contextTag,
				Name:      key,
			}, tagKeys,
		)
	}
	return &gauge{
		collector: newCollector(registerer, createGaugeFn),
	}
}

// Set gauge with the value.
func (g *gauge) Set(contextTag, key string, value float64, tags ...Tag) error {
	gauge, err := g.getOrCreateGauge(contextTag, key, tags...)
	if err != nil {
		return err
	}
	gauge.Set(value)
	return nil
}

// Inc increases gauge by 1.
func (g *gauge) Inc(contextTag, key string, tags ...Tag) error {
	gauge, err := g.getOrCreateGauge(contextTag, key, tags...)
	if err != nil {
		return err
	}
	gauge.Inc()
	return nil
}

// Dec decreases gauge by 1.
func (g *gauge) Dec(contextTag, key string, tags ...Tag) error {
	gauge, err := g.getOrCreateGauge(contextTag, key, tags...)
	if err != nil {
		return err
	}
	gauge.Dec()
	return nil
}

func (g *gauge) getOrCreateGauge(contextTag, key string, tags ...Tag) (prometheus.Gauge, error) {
	collector := g.collector.GetOrCreateCollector(contextTag, key, tags...)
	gaugeVec := collector.(*prometheus.GaugeVec)

	tagMap := getTagMap(tags...)
	return gaugeVec.GetMetricWith(tagMap)
}
