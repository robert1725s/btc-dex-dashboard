package database

import (
	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/gorm"
)

// Seed は初期データを投入する
func Seed(db *gorm.DB) error {
	exchanges := []model.Exchange{
		{Key: "hyperliquid", DisplayName: "Hyperliquid"},
		{Key: "lighter", DisplayName: "Lighter"},
		{Key: "aster", DisplayName: "Aster"},
	}

	for _, ex := range exchanges {
		result := db.Where("key = ?", ex.Key).FirstOrCreate(&ex)
		if result.Error != nil {
			return result.Error
		}
	}

	var hyperliquid, lighter, aster model.Exchange
	db.Where("key = ?", "hyperliquid").First(&hyperliquid)
	db.Where("key = ?", "lighter").First(&lighter)
	db.Where("key = ?", "aster").First(&aster)

	markets := []model.Market{
		{ExchangeID: hyperliquid.ID, Symbol: "BTC", BaseAsset: "BTC", QuoteAsset: "USDT"},
		{ExchangeID: lighter.ID, Symbol: "BTC-PERP", BaseAsset: "BTC", QuoteAsset: "USDT"},
		{ExchangeID: aster.ID, Symbol: "BTCUSDT", BaseAsset: "BTC", QuoteAsset: "USDT"},
	}

	for _, m := range markets {
		result := db.Where("exchange_id = ? AND symbol = ?", m.ExchangeID, m.Symbol).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
