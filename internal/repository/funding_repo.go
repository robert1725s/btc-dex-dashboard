package repository

import (
	"context"
	"time"

	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/gorm"
)

type FundingRateRepository interface {
	FindByMarketAndTimeRange(ctx context.Context, marketID uint, from, to time.Time) ([]model.FundingRate, error)
	FindLatestByMarket(ctx context.Context, marketID uint) (*model.FundingRate, error)
	Create(ctx context.Context, rate *model.FundingRate) error
}

type GormFundingRateRepository struct {
	db *gorm.DB
}

func NewGormFundingRateRepository(db *gorm.DB) *GormFundingRateRepository {
	return &GormFundingRateRepository{db: db}
}

func (r *GormFundingRateRepository) FindByMarketAndTimeRange(ctx context.Context, marketID uint, from, to time.Time) ([]model.FundingRate, error) {
	var rates []model.FundingRate
	result := r.db.WithContext(ctx).
		Where("market_id = ? AND ts >= ? AND ts <= ?", marketID, from, to).
		Order("ts ASC").
		Find(&rates)
	return rates, result.Error
}

func (r *GormFundingRateRepository) FindLatestByMarket(ctx context.Context, marketID uint) (*model.FundingRate, error) {
	var rate model.FundingRate
	result := r.db.WithContext(ctx).
		Where("market_id = ?", marketID).
		Order("ts DESC").
		First(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rate, nil
}

func (r *GormFundingRateRepository) Create(ctx context.Context, rate *model.FundingRate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}