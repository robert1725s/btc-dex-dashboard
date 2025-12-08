package repository

import (
	"context"

	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/gorm"
)

// ExchangeRepository は取引所データへのアクセスを抽象化する interface
type ExchangeRepository interface {
	FindAll(ctx context.Context) ([]model.Exchange, error)
	FindByKey(ctx context.Context, key string) (*model.Exchange, error)
	Create(ctx context.Context, exchange *model.Exchange) error
}

// GormExchangeRepository は GORM を使った実装
type GormExchangeRepository struct {
	db *gorm.DB
}

// NewGormExchangeRepository はコンストラクタ
func NewGormExchangeRepository(db *gorm.DB) *GormExchangeRepository {
	return &GormExchangeRepository{db: db}
}

func (r *GormExchangeRepository) FindAll(ctx context.Context) ([]model.Exchange, error) {
	var exchanges []model.Exchange
	result := r.db.WithContext(ctx).Find(&exchanges)
	return exchanges, result.Error
}

func (r *GormExchangeRepository) FindByKey(ctx context.Context, key string) (*model.Exchange, error) {
	var exchange model.Exchange
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&exchange)
	if result.Error != nil {
		return nil, result.Error
	}
	return &exchange, nil
}

func (r *GormExchangeRepository) Create(ctx context.Context, exchange *model.Exchange) error {
	return r.db.WithContext(ctx).Create(exchange).Error
}
