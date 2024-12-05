//go:build wireinject
// +build wireinject

package mongo

import (
	"context"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeMongo(ctx context.Context, dns MongoConnectionString) *mongo.Client {
	wire.Build(newMongo)
	return nil
}
