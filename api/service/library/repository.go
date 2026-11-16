package library

import (
	"context"
	"errors"
	"geekai/store/model"

	"gorm.io/gorm"
)

type Repository interface {
	Sources(context.Context) ([]model.LibrarySource, error)
	Source(context.Context, uint) (model.LibrarySource, error)
	SaveSource(context.Context, *model.LibrarySource) error
	DeleteSource(context.Context, uint) error
	Resources(context.Context) ([]model.LibraryResource, error)
	Resource(context.Context, uint) (model.LibraryResource, error)
	SaveResource(context.Context, *model.LibraryResource) error
	DeleteResource(context.Context, uint) error
}

type GormRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &GormRepository{db: db} }

func recordError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func (r *GormRepository) Sources(ctx context.Context) ([]model.LibrarySource, error) {
	items := []model.LibrarySource{}
	err := r.db.WithContext(ctx).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *GormRepository) Source(ctx context.Context, id uint) (model.LibrarySource, error) {
	var item model.LibrarySource
	err := r.db.WithContext(ctx).First(&item, id).Error
	return item, recordError(err)
}

func (r *GormRepository) SaveSource(ctx context.Context, item *model.LibrarySource) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormRepository) DeleteSource(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.LibraryResource{}).Where("source_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return errors.New("请先删除该来源关联的发布记录")
		}
		return tx.Delete(&model.LibrarySource{}, id).Error
	})
}

func (r *GormRepository) Resources(ctx context.Context) ([]model.LibraryResource, error) {
	items := []model.LibraryResource{}
	err := r.db.WithContext(ctx).Order("sort_num ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *GormRepository) Resource(ctx context.Context, id uint) (model.LibraryResource, error) {
	var item model.LibraryResource
	err := r.db.WithContext(ctx).First(&item, id).Error
	return item, recordError(err)
}

func (r *GormRepository) SaveResource(ctx context.Context, item *model.LibraryResource) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormRepository) DeleteResource(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LibraryResource{}, id).Error
}
