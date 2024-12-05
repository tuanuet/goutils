package stats

import (
	"context"
	"io"

	"net/http"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"

	"github.com/stretchr/testify/assert"
)

type fakeHttpClient struct {
	requests []*http.Request
}

func (f *fakeHttpClient) Do(request *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, request)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("OK")),
	}, nil
}

func TestPusherClient(t *testing.T) {
	tests := []struct {
		name         string
		operation    func(c *prometheusPusherClient)
		expectedBody []string
	}{
		{
			name: "Test CountInc func without contextTag",
			operation: func(c *prometheusPusherClient) {
				c.CountInc(context.Background(), "", "count_1")
			},
			expectedBody: []string{"count_1 1"},
		},
		{
			name: "Test Count func with tags",
			operation: func(c *prometheusPusherClient) {
				c.Count(context.Background(), "svc", "count_1", 10, Tag{Key: "tag_key", Value: "value"})
			},
			expectedBody: []string{"svc_count_1{tag_key=\"value\"} 10"},
		},
		{
			name: "Test Gauge func",
			operation: func(c *prometheusPusherClient) {
				c.Gauge(context.Background(), "svc", "gauge_1", 1.2)
			},
			expectedBody: []string{"svc_gauge_1 1.2"},
		},
		{
			name: "Test GaugeInc func",
			operation: func(c *prometheusPusherClient) {
				c.GaugeInc(context.Background(), "svc", "gauge_1")
			},
			expectedBody: []string{"svc_gauge_1 1"},
		},
		{
			name: "Test GaugeDec func",
			operation: func(c *prometheusPusherClient) {
				c.GaugeDec(context.Background(), "svc", "gauge_1")
			},
			expectedBody: []string{"svc_gauge_1 -1"},
		},
		{
			name: "Test Histogram func",
			operation: func(c *prometheusPusherClient) {
				c.Histogram(context.Background(), "svc", "histogram_1", 11.5)
				c.Histogram(context.Background(), "svc", "histogram_1", 13)
			},
			expectedBody: []string{"svc_histogram_1_sum 24.5", "svc_histogram_1_count 2"},
		},
		{
			name: "Test multiple func",
			operation: func(c *prometheusPusherClient) {
				c.Count(context.Background(), "svc", "count_1", 9)
				c.Histogram(context.Background(), "", "histogram_1", 13)
				c.Gauge(context.Background(), "svc", "gauge_1", 8.9)
			},
			expectedBody: []string{
				"svc_count_1 9",
				"histogram_1_sum 13",
				"histogram_1_count 1",
				"svc_gauge_1 8.9",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeHttpClient := &fakeHttpClient{}

			cfg := PrometheusPusherConfig{
				Url:           "example.com",
				Job:           "job_name",
				DefaultLabels: map[string]string{"key1": "value"},
				Username:      "username",
				Password:      "password",
				Client:        fakeHttpClient,
				Format:        expfmt.FmtText,
			}
			pusherClient := NewPusherClient(cfg).(*prometheusPusherClient)

			tt.operation(pusherClient)

			lastRequest := fakeHttpClient.requests[len(fakeHttpClient.requests)-1]
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, lastRequest.Body)
			requestBody := buf.String()

			for _, expectString := range tt.expectedBody {
				assert.Contains(t, requestBody, expectString)
			}

		})
	}
}
