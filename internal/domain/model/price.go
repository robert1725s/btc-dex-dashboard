package model

import "time"

// Price1m は1分足の価格データ
type Price1m struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MarketID  uint      `gorm:"not null;uniqueIndex:idx_price_market_ts" json:"market_id"`
	Market    Market    `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	Ts        time.Time `gorm:"not null;uniqueIndex:idx_price_market_ts" json:"ts"`
	Close     float64   `gorm:"type:decimal(20,8);not null" json:"close"`
	CreatedAt time.Time `json:"created_at"`
}
