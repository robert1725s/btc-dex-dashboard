package model

import "time"

// FundingRate はFunding Rate（資金調達率）の履歴
type FundingRate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MarketID  uint      `gorm:"not null;uniqueIndex:idx_funding_market_ts" json:"market_id"`
	Market    Market    `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	Ts        time.Time `gorm:"not null;uniqueIndex:idx_funding_market_ts" json:"ts"`
	Rate      float64   `gorm:"type:decimal(20,10);not null" json:"rate"`
	CreatedAt time.Time `json:"created_at"`
}
