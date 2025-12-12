package service

import (
	"context"
	"sort"
	"time"

	"btc-dex-dashboard/internal/repository"
)

type SpreadResult struct {
	Prices          []PriceInfo    `json:"prices"`
	BuyOpportunity  *ArbitrageInfo `json:"buy_opportunity"`
	SellOpportunity *ArbitrageInfo `json:"sell_opportunity"`
	History         []HistoryPoint `json:"history"`
	Stats           *SpreadStats   `json:"stats"`
}

type PriceInfo struct {
	ExchangeKey  string  `json:"exchange_key"`
	ExchangeName string  `json:"exchange_name"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	MidPrice     float64 `json:"mid_price"`
}

type HistoryPoint struct {
	Timestamp   string  `json:"timestamp"`
	Hyperliquid float64 `json:"hyperliquid"`
	Lighter     float64 `json:"lighter"`
	Aster       float64 `json:"aster"`
	Spread      float64 `json:"spread"`
	SpreadPct   float64 `json:"spread_pct"`
}

type MaxSpreadInfo struct {
	Value        float64 `json:"value"`
	Pct          float64 `json:"pct"`
	Timestamp    string  `json:"timestamp"`
	HighExchange string  `json:"high_exchange"`
	LowExchange  string  `json:"low_exchange"`
	HighPrice    float64 `json:"high_price"`
	LowPrice     float64 `json:"low_price"`
}

type SpreadStats struct {
	MaxSpread     *MaxSpreadInfo `json:"max_spread"`
	AvgSpread     float64        `json:"avg_spread"`
	AvgSpreadPct  float64        `json:"avg_spread_pct"`
	AvgPrice      float64        `json:"avg_price"`
	PeriodMinutes int            `json:"period_minutes"`
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

	// マーケットIDとキーのマッピング
	marketKeyToID := make(map[string]uint)

	for _, market := range markets {
		latestPrice, err := s.priceRepo.FindLatestByMarket(ctx, market.ID)
		if err != nil {
			continue
		}

		midPrice := (latestPrice.Bid + latestPrice.Ask) / 2
		info := PriceInfo{
			ExchangeKey:  market.Exchange.Key,
			ExchangeName: market.Exchange.DisplayName,
			Bid:          latestPrice.Bid,
			Ask:          latestPrice.Ask,
			MidPrice:     midPrice,
		}
		prices = append(prices, info)

		exchangePrices = append(exchangePrices, exchangePrice{
			name: market.Exchange.DisplayName,
			bid:  latestPrice.Bid,
			ask:  latestPrice.Ask,
		})

		marketKeyToID[market.Exchange.Key] = market.ID
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

	// 履歴データと統計情報を取得
	history, stats := s.calculateHistoryAndStats(ctx, marketKeyToID)

	return &SpreadResult{
		Prices:          prices,
		BuyOpportunity:  buyOpp,
		SellOpportunity: sellOpp,
		History:         history,
		Stats:           stats,
	}, nil
}

func (s *SpreadService) calculateHistoryAndStats(ctx context.Context, marketKeyToID map[string]uint) ([]HistoryPoint, *SpreadStats) {
	periodMinutes := 15
	now := time.Now()
	from := now.Add(-time.Duration(periodMinutes) * time.Minute)

	// 各取引所の履歴を取得
	type priceData struct {
		ts       time.Time
		midPrice float64
	}
	exchangeHistory := make(map[string][]priceData)

	for key, marketID := range marketKeyToID {
		prices, err := s.priceRepo.FindByMarketAndTimeRange(ctx, marketID, from, now)
		if err != nil {
			continue
		}
		for _, p := range prices {
			exchangeHistory[key] = append(exchangeHistory[key], priceData{
				ts:       p.Ts,
				midPrice: (p.Bid + p.Ask) / 2,
			})
		}
	}

	// タイムスタンプでマージして履歴ポイントを作成
	// 全タイムスタンプを収集
	allTimestamps := make(map[time.Time]bool)
	for _, prices := range exchangeHistory {
		for _, p := range prices {
			// 秒単位で丸める
			ts := p.ts.Truncate(time.Second)
			allTimestamps[ts] = true
		}
	}

	// ソート
	var sortedTimestamps []time.Time
	for ts := range allTimestamps {
		sortedTimestamps = append(sortedTimestamps, ts)
	}
	sort.Slice(sortedTimestamps, func(i, j int) bool {
		return sortedTimestamps[i].Before(sortedTimestamps[j])
	})

	// 各取引所の価格を時刻でインデックス化
	exchangePriceByTime := make(map[string]map[time.Time]float64)
	for key, prices := range exchangeHistory {
		exchangePriceByTime[key] = make(map[time.Time]float64)
		for _, p := range prices {
			ts := p.ts.Truncate(time.Second)
			exchangePriceByTime[key][ts] = p.midPrice
		}
	}

	// 履歴ポイントを生成
	var history []HistoryPoint
	var maxSpread *MaxSpreadInfo
	var totalSpread, totalPrice float64
	var spreadCount, priceCount int

	// 直前の価格を保持（欠損値対応）
	lastPrices := make(map[string]float64)

	for _, ts := range sortedTimestamps {
		point := HistoryPoint{
			Timestamp: ts.Format(time.RFC3339),
		}

		// 各取引所の価格を設定（なければ直前の値を使用）
		for _, key := range []string{"hyperliquid", "lighter", "aster"} {
			if price, ok := exchangePriceByTime[key][ts]; ok {
				lastPrices[key] = price
			}
		}

		point.Hyperliquid = lastPrices["hyperliquid"]
		point.Lighter = lastPrices["lighter"]
		point.Aster = lastPrices["aster"]

		// スプレッド計算（最高値 - 最低値）
		var validPrices []struct {
			name  string
			price float64
		}
		if point.Hyperliquid > 0 {
			validPrices = append(validPrices, struct {
				name  string
				price float64
			}{"HyperLiquid", point.Hyperliquid})
		}
		if point.Lighter > 0 {
			validPrices = append(validPrices, struct {
				name  string
				price float64
			}{"Lighter", point.Lighter})
		}
		if point.Aster > 0 {
			validPrices = append(validPrices, struct {
				name  string
				price float64
			}{"Aster", point.Aster})
		}

		if len(validPrices) >= 2 {
			// ソートして最高値と最低値を取得
			sort.Slice(validPrices, func(i, j int) bool {
				return validPrices[i].price > validPrices[j].price
			})

			highPrice := validPrices[0].price
			lowPrice := validPrices[len(validPrices)-1].price
			spread := highPrice - lowPrice
			spreadPct := (spread / lowPrice) * 100

			point.Spread = spread
			point.SpreadPct = spreadPct

			totalSpread += spread
			spreadCount++

			// 平均価格の計算用
			for _, vp := range validPrices {
				totalPrice += vp.price
				priceCount++
			}

			// 最大スプレッド更新
			if maxSpread == nil || spread > maxSpread.Value {
				maxSpread = &MaxSpreadInfo{
					Value:        spread,
					Pct:          spreadPct,
					Timestamp:    ts.Format(time.RFC3339),
					HighExchange: validPrices[0].name,
					LowExchange:  validPrices[len(validPrices)-1].name,
					HighPrice:    highPrice,
					LowPrice:     lowPrice,
				}
			}
		}

		history = append(history, point)
	}

	// 3つの取引所すべての価格があるポイントのみをフィルタ
	var validHistory []HistoryPoint
	for _, p := range history {
		if p.Hyperliquid > 0 && p.Lighter > 0 && p.Aster > 0 {
			validHistory = append(validHistory, p)
		}
	}
	history = validHistory

	// サンプリング（最大180ポイント）
	maxPoints := 180
	if len(history) > maxPoints {
		step := len(history) / maxPoints
		var sampled []HistoryPoint
		for i := 0; i < len(history); i += step {
			sampled = append(sampled, history[i])
		}
		history = sampled
	}

	// 統計情報を計算
	var avgSpread, avgSpreadPct, avgPrice float64
	if spreadCount > 0 {
		avgSpread = totalSpread / float64(spreadCount)
	}
	if priceCount > 0 {
		avgPrice = totalPrice / float64(priceCount)
	}
	if avgPrice > 0 {
		avgSpreadPct = (avgSpread / avgPrice) * 100
	}

	stats := &SpreadStats{
		MaxSpread:     maxSpread,
		AvgSpread:     avgSpread,
		AvgSpreadPct:  avgSpreadPct,
		AvgPrice:      avgPrice,
		PeriodMinutes: periodMinutes,
	}

	return history, stats
}