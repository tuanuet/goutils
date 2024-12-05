package stats

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// collector is used to bundle many prometheus.Collector instances with differ in contextTag, key and tagKeys but have
// same type. It provides a thread safe mechanism to get and create new collector. collector is not used directly. It is
// used as a building block for implementations of counter, histogram and gauge
type collector struct {
	registerer        prometheus.Registerer
	collectorMap      map[string]prometheus.Collector
	mux               sync.Mutex
	createCollectorFn func(contextTag, key string, tagKeys ...string) prometheus.Collector
}

func newCollector(
	registerer prometheus.Registerer,
	createCollectorFn func(contextTag, key string, tagKeys ...string) prometheus.Collector) *collector {
	return &collector{
		registerer:        registerer,
		createCollectorFn: createCollectorFn,
		collectorMap:      make(map[string]prometheus.Collector),
	}
}

// GetOrCreateCollector is thread safe.
func (c *collector) GetOrCreateCollector(contextTag, key string, tags ...Tag) prometheus.Collector {
	sortedTagKeys := getSortedTagKeys(tags...)
	hashKey := buildHashKey(contextTag, key, sortedTagKeys...)

	collector, isExist := c.collectorMap[hashKey]

	if !isExist {
		c.mux.Lock()
		collector, isExist = c.collectorMap[hashKey]
		if !isExist {
			collector = c.createCollectorFn(contextTag, key, sortedTagKeys...)
			c.addCollector(hashKey, collector)
		}
		c.mux.Unlock()
	}

	return collector
}

func (c *collector) addCollector(hashKey string, collector prometheus.Collector) {
	c.collectorMap[hashKey] = collector
	c.registerer.MustRegister(collector)
}
