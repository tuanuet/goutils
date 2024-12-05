package postgres

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type PostgresConnectionString string

func newPostgresGorm(ctx context.Context, dns PostgresConnectionString) *gorm.DB {
	var once sync.Once
	once.Do(func() {
		_db, err := gorm.Open(postgres.Open(string(dns)), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("cannot connect database: %+v", err))
		}
		db = _db
	})
	return db
}
