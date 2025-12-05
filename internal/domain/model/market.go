package model

import "time"

// Market はマーケット（取引ペア）
type Market struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ExchangeID uint      `gorm:"not null;index" json:"exchange_id"`
	Exchange   Exchange  `gorm:"foreignKey:ExchangeID" json:"exchange,omitempty"`
	Symbol     string    `gorm:"size:50;not null" json:"symbol"`
	BaseAsset  string    `gorm:"size:20;not null" json:"base_asset"`
	QuoteAsset string    `gorm:"size:20;not null" json:"quote_asset"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
