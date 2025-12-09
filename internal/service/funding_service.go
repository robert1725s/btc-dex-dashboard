package service

import (
	"context"

	"btc-dex-dashboard/internal/repository"
)

type FundingResult struct {
	Rates []FundingInfo `json:"rates"`
}

type FundingInfo struct {
	ExchangeKey  string  `json:"exchange_key"`
	ExchangeName string  `json:"exchange_name"`
	Rate         float64 `json:"rate"`
	RatePct      float64 `json:"rate_pct"`
}

type FundingService struct {
	marketRepo  repository.MarketRepository
	fundingRepo repository.FundingRateRepository
}

func NewFundingService(
	marketRepo repository.MarketRepository,
	fundingRepo repository.FundingRateRepository,
) *FundingService {
	return &FundingService{
		marketRepo:  marketRepo,
		fundingRepo: fundingRepo,
	}
}

func (s *FundingService) GetLatestRates(ctx context.Context) (*FundingResult, error) {
	markets, err := s.marketRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var rates []FundingInfo

	for _, market := range markets {
		latestRate, err := s.fundingRepo.FindLatestByMarket(ctx, market.ID)
		if err != nil {
			continue
		}

		info := FundingInfo{
			ExchangeKey:  market.Exchange.Key,
			ExchangeName: market.Exchange.DisplayName,
			Rate:         latestRate.Rate,
			RatePct:      latestRate.Rate * 100,
		}
		rates = append(rates, info)
	}

	return &FundingResult{Rates: rates}, nil
}