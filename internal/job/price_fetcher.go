package job

import (
	"context"
	"log"
	"sync"

	"btc-dex-dashboard/internal/domain/model"
	"btc-dex-dashboard/internal/infrastructure/dex"
	"btc-dex-dashboard/internal/repository"
)

type PriceFetcher struct {
	clients   []dex.DexClient
	priceRepo repository.PriceRepository
	marketIDs map[string]uint // DEX名 → MarketID のマッピング
}

func NewPriceFetcher(
	clients []dex.DexClient,
	priceRepo repository.PriceRepository,
	marketIDs map[string]uint,
) *PriceFetcher {
	return &PriceFetcher{
		clients:   clients,
		priceRepo: priceRepo,
		marketIDs: marketIDs,
	}
}

type priceResult struct {
	dexName string
	data    *dex.PriceData
	err     error
}

func (f *PriceFetcher) FetchAndSaveAll(ctx context.Context) {
	results := make(chan priceResult, len(f.clients))
	var wg sync.WaitGroup

	// 全 DEX から並行して価格を取得
	for _, client := range f.clients {
		wg.Add(1)
		go func(c dex.DexClient) {
			defer wg.Done()

			data, err := c.FetchBTCPerpPrice(ctx)
			results <- priceResult{
				dexName: c.Name(),
				data:    data,
				err:     err,
			}
		}(client)
	}

	// 全 goroutine の完了を待ってから channel を閉じる
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を受信して DB に保存
	for result := range results {
		if result.err != nil {
			log.Printf("[%s] failed to fetch price: %v", result.dexName, result.err)
			continue
		}

		marketID, ok := f.marketIDs[result.dexName]
		if !ok {
			log.Printf("[%s] market ID not found", result.dexName)
			continue
		}

		price := &model.Price{
			MarketID: marketID,
			Ts:       result.data.Ts,
			Bid:      result.data.Bid,
			Ask:      result.data.Ask,
		}

		if err := f.priceRepo.Create(ctx, price); err != nil {
			log.Printf("[%s] failed to save price: %v", result.dexName, err)
			continue
		}

		log.Printf("[%s] saved: bid=%.2f, ask=%.2f", result.dexName, result.data.Bid, result.data.Ask)
	}
}
