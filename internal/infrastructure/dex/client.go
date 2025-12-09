package dex

import (
	"context"
	"time"
)

type PriceData struct {
	Bid float64
	Ask float64
	Ts  time.Time
}

type FundingRateData struct {
	Rate float64
	Ts   time.Time
}

type DexClient interface {
	Name() string
	FetchBTCPerpPrice(ctx context.Context) (*PriceData, error)
	FetchBTCPerpFundingRate(ctx context.Context) (*FundingRateData, error)
}
