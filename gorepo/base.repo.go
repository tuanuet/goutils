package gorepo

import "gorm.io/gorm/schema"

type BaseRepo[E schema.Tabler, ID any] interface {
	// FindById ... find one by ID, return entity, isExisted, error
	FindById(id ID) (*E, bool, error)

	// FindByIDs ...
	FindByIDs(ids []ID) ([]E, error)

	// FindAll ... find all data
	FindAll() ([]E, error)

	// FindByColumn ... find in {column} that match {value}
	FindByColumn(colName string, value string) ([]E, error)

	// FindByColumnIn ... find in {column} that match {value}
	FindByColumnIn(colName string, values []string) ([]E, error)

	// Create ... create new entity
	Create(entity *E) error

	// Upsert ... insert or update entity
	Upsert(entity *E) error

	// UpsertAll ... create new entities by batch
	UpsertAll(entities []*E) error

	// UpsertBatch ... create new entities by batch
	UpsertBatch(entities []*E, batchSize int) error

	// DeleteAll ... delete by list ID
	DeleteAll(entities []E) error

	// CountAll ... count all data
	CountAll() (int, error)
}
