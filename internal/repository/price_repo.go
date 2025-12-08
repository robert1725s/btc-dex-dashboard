package repository

import (
	"context"
	"time"

	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/gorm"
)

type PriceRepository interface {
	FindByMarketAndTimeRange(ctx context.Context, marketID uint, from, to time.Time) ([]model.Price, error)
	FindLatestByMarket(ctx context.Context, marketID uint) (*model.Price, error)
	Create(ctx context.Context, price *model.Price) error
	CreateBatch(ctx context.Context, prices []model.Price) error
}

type GormPriceRepository struct {
	db *gorm.DB
}

func NewGormPriceRepository(db *gorm.DB) *GormPriceRepository {
	return &GormPriceRepository{db: db}
}

func (r *GormPriceRepository) FindByMarketAndTimeRange(ctx context.Context, marketID uint, from, to time.Time) ([]model.Price, error) {
	var prices []model.Price
	result := r.db.WithContext(ctx).
		Where("market_id = ? AND ts >= ? AND ts <= ?", marketID, from, to).
		Order("ts ASC").
		Find(&prices)
	return prices, result.Error
}

func (r *GormPriceRepository) FindLatestByMarket(ctx context.Context, marketID uint) (*model.Price, error) {
	var price model.Price
	result := r.db.WithContext(ctx).
		Where("market_id = ?", marketID).
		Order("ts DESC").
		First(&price)
	if result.Error != nil {
		return nil, result.Error
	}
	return &price, nil
}

func (r *GormPriceRepository) Create(ctx context.Context, price *model.Price) error {
	return r.db.WithContext(ctx).Create(price).Error
}

func (r *GormPriceRepository) CreateBatch(ctx context.Context, prices []model.Price) error {
	return r.db.WithContext(ctx).Create(&prices).Error
}
