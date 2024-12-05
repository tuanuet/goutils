//go:build wireinject
// +build wireinject

package postgres

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializePostgres(ctx context.Context, dns PostgresConnectionString) *gorm.DB {
	wire.Build(newPostgresGorm)
	return nil
}
