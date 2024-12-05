// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package stats

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, contextTag, key, amount, tags
func (_m *MockClient) Count(ctx context.Context, contextTag string, key string, amount int64, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key, amount)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// CountInc provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *MockClient) CountInc(ctx context.Context, contextTag string, key string, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Duration provides a mock function with given fields: ctx, contextTag, key, start, tags
func (_m *MockClient) Duration(ctx context.Context, contextTag string, key string, start time.Time, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key, start)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Gauge provides a mock function with given fields: ctx, contextTag, key, value, tags
func (_m *MockClient) Gauge(ctx context.Context, contextTag string, key string, value float64, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key, value)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// GaugeDec provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *MockClient) GaugeDec(ctx context.Context, contextTag string, key string, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// GaugeInc provides a mock function with given fields: ctx, contextTag, key, tags
func (_m *MockClient) GaugeInc(ctx context.Context, contextTag string, key string, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Histogram provides a mock function with given fields: ctx, contextTag, key, value, tags
func (_m *MockClient) Histogram(ctx context.Context, contextTag string, key string, value float64, tags ...Tag) {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contextTag, key, value)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}
