package repository

import (
	"context"

	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/gorm"
)

type MarketRepository interface {
	FindAll(ctx context.Context) ([]model.Market, error)
	FindByExchangeID(ctx context.Context, exchangeID uint) ([]model.Market, error)
	FindByID(ctx context.Context, id uint) (*model.Market, error)
	Create(ctx context.Context, market *model.Market) error
}

type GormMarketRepository struct {
	db *gorm.DB
}

func NewGormMarketRepository(db *gorm.DB) *GormMarketRepository {
	return &GormMarketRepository{db: db}
}

func (r *GormMarketRepository) FindAll(ctx context.Context) ([]model.Market, error) {
	var markets []model.Market
	result := r.db.WithContext(ctx).Preload("Exchange").Find(&markets)
	return markets, result.Error
}

func (r *GormMarketRepository) FindByExchangeID(ctx context.Context, exchangeID uint) ([]model.Market, error) {
	var markets []model.Market
	result := r.db.WithContext(ctx).Where("exchange_id = ?", exchangeID).Find(&markets)
	return markets, result.Error
}

func (r *GormMarketRepository) FindByID(ctx context.Context, id uint) (*model.Market, error) {
	var market model.Market
	result := r.db.WithContext(ctx).Preload("Exchange").First(&market, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &market, nil
}

func (r *GormMarketRepository) Create(ctx context.Context, market *model.Market) error {
	return r.db.WithContext(ctx).Create(market).Error
}
