package clickhouse

import (
	"context"
	"fmt"
	clickhousedriver "github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

var (
	db *gorm.DB
)

type ClickhouseConnectionString string

type ConnectionOpt struct {
	Addr     string
	DBName   string
	Username string
	Password string
}

func newClickhouseGorm(ctx context.Context, opt ConnectionOpt) *gorm.DB {
	var once sync.Once
	once.Do(func() {
		sqlDB := clickhousedriver.OpenDB(&clickhousedriver.Options{
			Addr: strings.Split(opt.Addr, ","),
			Auth: clickhousedriver.Auth{
				Database: opt.DBName,
				Username: opt.Username,
				Password: opt.Password,
			},
			ConnOpenStrategy: clickhousedriver.ConnOpenRoundRobin,
		})

		sqlDB.SetConnMaxLifetime(time.Duration(30) * time.Minute)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)

		if err := sqlDB.Ping(); err != nil {
			panic(err)
		}

		_db, err := gorm.Open(clickhouse.New(clickhouse.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("cannot connect database: %+v", err))
		}
		db = _db
	})
	return db
}
