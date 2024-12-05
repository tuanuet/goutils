package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Client

type MongoConnectionString string

func newMongo(ctx context.Context, dns MongoConnectionString) *mongo.Client {
	var once sync.Once
	once.Do(func() {
		_db, err := mongo.Connect(ctx, options.Client().ApplyURI(string(dns)))
		if err != nil {
			panic(fmt.Errorf("cannot connect mongo: %+v", err))
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err = _db.Ping(timeoutCtx, readpref.Nearest()); err != nil {
			panic(fmt.Errorf("cannot ping mongo: %+v", err))
		}

		fmt.Println("[Connected] mongo connected")
		db = _db
	})

	return db
}
