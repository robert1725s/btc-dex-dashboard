package database

import (
	"btc-dex-dashboard/internal/domain/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB は新しいデータベース接続を作成する
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.Exchange{},
		&model.Market{},
		&model.Price1m{},
		&model.FundingRate{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
