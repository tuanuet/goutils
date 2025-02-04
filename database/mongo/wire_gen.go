// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitializeMongo(ctx context.Context, dns MongoConnectionString) *mongo.Client {
	client := newMongo(ctx, dns)
	return client
}
