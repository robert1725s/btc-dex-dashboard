package service

import (
	"context"

	"btc-dex-dashboard/internal/repository"
)

type SpreadResult struct {
	Prices       []PriceInfo       `json:"prices"`
	BuyOpportunity  *ArbitrageInfo `json:"buy_opportunity"`
	SellOpportunity *ArbitrageInfo `json:"sell_opportunity"`
}

type PriceInfo struct {
	ExchangeKey  string  `json:"exchange_key"`
	ExchangeName string  `json:"exchange_name"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
}

type ArbitrageInfo struct {
	BuyExchange  string  `json:"buy_exchange"`
	SellExchange string  `json:"sell_exchange"`
	BuyPrice     float64 `json:"buy_price"`
	SellPrice    float64 `json:"sell_price"`
	SpreadAbs    float64 `json:"spread_abs"`
	SpreadPct    float64 `json:"spread_pct"`
}

type SpreadService struct {
	marketRepo repository.MarketRepository
	priceRepo  repository.PriceRepository
}

func NewSpreadService(
	marketRepo repository.MarketRepository,
	priceRepo repository.PriceRepository,
) *SpreadService {
	return &SpreadService{
		marketRepo: marketRepo,
		priceRepo:  priceRepo,
	}
}

func (s *SpreadService) CalculateSpread(ctx context.Context) (*SpreadResult, error) {
	markets, err := s.marketRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var prices []PriceInfo

	type exchangePrice struct {
		name string
		bid  float64
		ask  float64
	}
	var exchangePrices []exchangePrice

	for _, market := range markets {
		latestPrice, err := s.priceRepo.FindLatestByMarket(ctx, market.ID)
		if err != nil {
			continue
		}

		info := PriceInfo{
			ExchangeKey:  market.Exchange.Key,
			ExchangeName: market.Exchange.DisplayName,
			Bid:          latestPrice.Bid,
			Ask:          latestPrice.Ask,
		}
		prices = append(prices, info)

		exchangePrices = append(exchangePrices, exchangePrice{
			name: market.Exchange.DisplayName,
			bid:  latestPrice.Bid,
			ask:  latestPrice.Ask,
		})
	}

	var buyOpp, sellOpp *ArbitrageInfo

	for i, ep1 := range exchangePrices {
		for j, ep2 := range exchangePrices {
			if i == j {
				continue
			}

			// ep1で買って(ask)、ep2で売る(bid)
			spread := ep2.bid - ep1.ask
			if spread > 0 {
				pct := (spread / ep1.ask) * 100
				if buyOpp == nil || spread > buyOpp.SpreadAbs {
					buyOpp = &ArbitrageInfo{
						BuyExchange:  ep1.name,
						SellExchange: ep2.name,
						BuyPrice:     ep1.ask,
						SellPrice:    ep2.bid,
						SpreadAbs:    spread,
						SpreadPct:    pct,
					}
				}
			}

			// ep1で売って(bid)、ep2で買う(ask) → 逆方向
			spreadRev := ep1.bid - ep2.ask
			if spreadRev > 0 {
				pct := (spreadRev / ep2.ask) * 100
				if sellOpp == nil || spreadRev > sellOpp.SpreadAbs {
					sellOpp = &ArbitrageInfo{
						BuyExchange:  ep2.name,
						SellExchange: ep1.name,
						BuyPrice:     ep2.ask,
						SellPrice:    ep1.bid,
						SpreadAbs:    spreadRev,
						SpreadPct:    pct,
					}
				}
			}
		}
	}

	return &SpreadResult{
		Prices:          prices,
		BuyOpportunity:  buyOpp,
		SellOpportunity: sellOpp,
	}, nil
}