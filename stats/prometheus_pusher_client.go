package stats

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
)

// HTTPDoer is an interface for the one method of http.Client that is used by Pusher.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// PrometheusPusherConfig contains params for creating new prometheusPusherClient.
type PrometheusPusherConfig struct {
	// required options
	Url string `help:"Promethues Pusher URL" env:"PROMETHEUS_PUSHER_URL" default:"localhost:9091"`
	Job string `help:"Promethues Pusher Job" env:"PROMETHEUS_PUSHER_JOB" default:"na"`

	// optional options
	DefaultLabels map[string]string `kong:"-"`
	Username      string            `kong:"-"`
	Password      string            `kong:"-"`
	Client        HTTPDoer          `kong:"-"`
	Format        expfmt.Format     `kong:"-"`
}

// prometheusPusherClient is an implementation of Client. It was created for the purpose of updating the goutils by
// pushing up to the Pushgateway service of prometheus
// https://prometheus.io/docs/practices/pushing/.
type prometheusPusherClient struct {
	client *prometheusClient
	pusher *push.Pusher
}

func NewPusherClient(cfg PrometheusPusherConfig) Client {
	registry := prometheus.NewRegistry()

	return &prometheusPusherClient{
		client: &prometheusClient{
			registerer: registry,
			counter:    newCounter(registry),
			histogram:  newHistogram(registry),
			gauge:      newGauge(registry),
		},
		pusher: newPusher(cfg).Gatherer(registry),
	}

}

func newPusher(cfg PrometheusPusherConfig) *push.Pusher {
	pusher := push.New(cfg.Url, cfg.Job)

	for label, value := range cfg.DefaultLabels {
		pusher.Grouping(label, value)
	}

	if len(cfg.Username) > 0 {
		pusher.BasicAuth(cfg.Username, cfg.Password)
	}

	if cfg.Client != nil {
		pusher.Client(cfg.Client)
	}

	if len(cfg.Format) > 0 {
		pusher.Format(cfg.Format)
	}

	return pusher
}

func (p *prometheusPusherClient) CountInc(ctx context.Context, contextTag, key string, tags ...Tag) {
	p.client.CountInc(ctx, contextTag, key, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) Count(ctx context.Context, contextTag, key string, amount int64, tags ...Tag) {
	p.client.Count(ctx, contextTag, key, amount, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) Histogram(ctx context.Context, contextTag, key string, value float64, tags ...Tag) {
	p.client.Histogram(ctx, contextTag, key, value, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) Duration(ctx context.Context, contextTag, key string, start time.Time, tags ...Tag) {
	p.client.Duration(ctx, contextTag, key, start, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) Gauge(ctx context.Context, contextTag, key string, value float64, tags ...Tag) {
	p.client.Gauge(ctx, contextTag, key, value, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) GaugeInc(ctx context.Context, contextTag, key string, tags ...Tag) {
	p.client.GaugeInc(ctx, contextTag, key, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) GaugeDec(ctx context.Context, contextTag, key string, tags ...Tag) {
	p.client.GaugeDec(ctx, contextTag, key, tags...)
	p.updateToPushGateway()
}

func (p *prometheusPusherClient) updateToPushGateway() {
	if err := p.pusher.Add(); err != nil {
		logrus.Warnf("cannot update metric on pushgateway server cause %v", err)
	}
}
