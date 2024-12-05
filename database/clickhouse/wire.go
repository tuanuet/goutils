//go:build wireinject
// +build wireinject

package clickhouse

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeClickhouse(ctx context.Context, opt ConnectionOpt) *gorm.DB {
	wire.Build(newClickhouseGorm)
	return nil
}
