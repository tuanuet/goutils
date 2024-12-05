package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildHashKey(t *testing.T) {
	type args struct {
		contextTag    string
		key           string
		sortedTagKeys []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Case",
			args: args{
				contextTag:    "svc",
				key:           "http_request_total",
				sortedTagKeys: []string{"appId", "url"},
			},
			want: "svc_http_request_total_appId_url",
		},
		{
			name: "contextTag is empty",
			args: args{
				contextTag:    "",
				key:           "http_request_total",
				sortedTagKeys: []string{"appId", "url"},
			},
			want: "http_request_total_appId_url",
		},
		{
			name: "sortedTagKeys is empty",
			args: args{
				contextTag:    "svc",
				key:           "http_request_total",
				sortedTagKeys: []string{},
			},
			want: "svc_http_request_total",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildHashKey(tt.args.contextTag, tt.args.key, tt.args.sortedTagKeys...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getSortedTagKeys(t *testing.T) {
	type args struct {
		tags []Tag
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Happy Case",
			args: args{
				tags: []Tag{
					{Key: "2key", Value: "value2"},
					{Key: "1key", Value: "value1"},
					{Key: "3key", Value: "value2"},
				},
			},
			want: []string{"1key", "2key", "3key"},
		},
		{
			name: "Empty tags",
			args: args{
				tags: []Tag{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSortedTagKeys(tt.args.tags...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getTagMap(t *testing.T) {
	type args struct {
		tags []Tag
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Happy Case",
			args: args{
				tags: []Tag{
					{Key: "appId", Value: "123"},
					{Key: "url", Value: "/test"},
				},
			},
			want: map[string]string{"appId": "123", "url": "/test"},
		},
		{
			name: "Missing tags",
			args: args{
				tags: []Tag{},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTagMap(tt.args.tags...)
			assert.Equal(t, tt.want, got)
		})
	}
}
