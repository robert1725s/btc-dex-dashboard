package model

import "time"

// Price は価格データ（リアルタイム / 秒単位）
type Price struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MarketID  uint      `gorm:"not null;uniqueIndex:idx_price_market_ts" json:"market_id"`
	Market    Market    `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	Ts        time.Time `gorm:"not null;uniqueIndex:idx_price_market_ts" json:"ts"`
	Price     float64   `gorm:"type:decimal(20,8);not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
}