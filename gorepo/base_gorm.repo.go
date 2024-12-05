package gorepo

import (
	"errors"
	"github.com/tuanuet/goutils/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// baseGormRepoImpl ... implement BaseRepo
type baseGormRepoImpl[E schema.Tabler, ID any] struct {
	db    *gorm.DB
	eName string // entity name
}

// NewBaseGormRepo ...
func NewBaseGormRepo[E schema.Tabler, ID any](db *gorm.DB, entityName string) BaseRepo[E, ID] {
	return &baseGormRepoImpl[E, ID]{db, entityName}
}

// Public methods ------------------------------------------------------------------------------------------------------

func (r *baseGormRepoImpl[E, ID]) FindById(id ID) (*E, bool, error) {
	var entity E
	if err := r.db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("Not found %v id=%v", r.eName, id)
			return nil, false, nil
		} else {
			logger.Errorf("Get %v got error %v", r.eName, err)
			return nil, false, err
		}
	}

	return &entity, true, nil
}

func (r *baseGormRepoImpl[E, ID]) FindByIDs(ids []ID) ([]E, error) {
	if len(ids) <= 0 {
		logger.Warnf("Get %v with ids is empty -> return empty", r.eName)
		return []E{}, nil
	}

	var entities []E
	if err := r.db.Find(&entities, ids).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("Not found %v id in (%v)", r.eName, ids)
			return []E{}, nil
		} else {
			logger.Errorf("Get %v got error %v", r.eName, err)
			return []E{}, nil
		}
	}

	return entities, nil
}

func (r *baseGormRepoImpl[E, ID]) FindAll() ([]E, error) {
	var entities []E
	if err := r.db.Find(&entities).Error; err != nil {
		logger.Errorf("Get all %v got error %v", r.eName, err)
		return nil, err
	}

	return entities, nil
}

func (r *baseGormRepoImpl[E, ID]) FindByColumn(colName string, value string) ([]E, error) {
	var entities []E
	if err := r.db.Where(colName+" = ?", value).Find(&entities).Error; err != nil {
		logger.Errorf("Get %v by %v = %v got error %v", r.eName, colName, value, err)
		return nil, err
	}

	return entities, nil
}

func (r *baseGormRepoImpl[E, ID]) FindByColumnIn(colName string, values []string) ([]E, error) {
	var entities []E
	if err := r.db.Where(colName+" IN (?)", values).Find(&entities).Error; err != nil {
		logger.Errorf("Get %v by %v in (%v) got error %v", r.eName, colName, values, err)
		return nil, err
	}

	return entities, nil
}

func (r *baseGormRepoImpl[E, ID]) Create(entity *E) error {
	if err := r.db.Save(entity).Error; err != nil {
		logger.Errorf("Create %v got error %s", r.eName, err.Error())
		return err
	}

	return nil
}

func (r *baseGormRepoImpl[E, ID]) Upsert(entity *E) error {
	err := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(entity).Error

	if err != nil {
		logger.Errorf("Upsert %v got error %s", r.eName, err.Error())
		return err
	}

	return nil
}

func (r *baseGormRepoImpl[E, ID]) UpsertAll(entities []*E) error {
	err := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(entities).Error

	if err != nil {
		logger.Errorf("Upsert list %v got error %s", r.eName, err.Error())
		return err
	}

	return nil
}

func (r *baseGormRepoImpl[E, ID]) UpsertBatch(entities []*E, batchSize int) error {
	return batchExecuting(batchSize, entities, r.UpsertAll)
}

func (r *baseGormRepoImpl[E, ID]) DeleteAll(entities []E) error {
	if err := r.db.Delete(&entities).Error; err != nil {
		logger.Errorf("Delete %v got error %s", r.eName, err.Error())
		return err
	}

	return nil
}

func (r *baseGormRepoImpl[E, ID]) CountAll() (int, error) {
	var count int64
	var entity E
	if err := r.db.Model(&entity).Count(&count).Error; err != nil {
		logger.Errorf("CountAll %v got error %s", r.eName, err.Error())
		return 0, err
	}

	return int(count), nil
}

func batchExecuting[T any](batchSize int, listData []T, execute func([]T) error) error {
	for head := 0; head < len(listData); head += batchSize {
		tail := head + batchSize
		if tail > len(listData) {
			tail = len(listData)
		}

		batch := listData[head:tail]

		logger.Debugf("Execute %v item (%v -> %v)\n", len(batch), head, tail)
		err := execute(batch)

		if err != nil {
			return err
		}

	}

	return nil
}
