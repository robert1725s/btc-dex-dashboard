package dex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const lighterAPIURL = "https://mainnet.zklighter.elliot.ai/api/v1/orderBookOrders"

type LighterClient struct {
	httpClient *http.Client
}

func NewLighterClient() *LighterClient {
	return &LighterClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *LighterClient) Name() string {
	return "lighter"
}

type lighterOrderBookResponse struct {
	Code      int            `json:"code"`
	TotalAsks int            `json:"total_asks"`
	Asks      []lighterOrder `json:"asks"`
	TotalBids int            `json:"total_bids"`
	Bids      []lighterOrder `json:"bids"`
}

type lighterOrder struct {
	OrderID         string `json:"order_id"`
	RemainingAmount string `json:"remaining_base_amount"`
	Price           string `json:"price"`
}

func (c *LighterClient) FetchBTCPerpPrice(ctx context.Context) (*PriceData, error) {
	// market_id=1 は BTC-PERP
	url := lighterAPIURL + "?market_id=1&limit=1"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var orderBookResp lighterOrderBookResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderBookResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(orderBookResp.Bids) == 0 || len(orderBookResp.Asks) == 0 {
		return nil, fmt.Errorf("invalid response: no bids or asks")
	}

	bid, err := strconv.ParseFloat(orderBookResp.Bids[0].Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid price: %w", err)
	}

	ask, err := strconv.ParseFloat(orderBookResp.Asks[0].Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask price: %w", err)
	}

	return &PriceData{
		Bid: bid,
		Ask: ask,
		Ts:  time.Now(),
	}, nil
}

func (c *LighterClient) FetchBTCPerpFundingRate(ctx context.Context) (*FundingRateData, error) {
	// TODO: Funding Rate 取得を実装
	return nil, fmt.Errorf("not implemented")
}
