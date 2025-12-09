package dex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const asterAPIURL = "https://fapi.asterdex.com/fapi/v1/ticker/bookTicker"

type AsterClient struct {
	httpClient *http.Client
}

func NewAsterClient() *AsterClient {
	return &AsterClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AsterClient) Name() string {
	return "aster"
}

type asterBookTickerResponse struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Time     int64  `json:"time"`
}

func (c *AsterClient) FetchBTCPerpPrice(ctx context.Context) (*PriceData, error) {
	url := asterAPIURL + "?symbol=BTCUSDT"

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

	var tickerResp asterBookTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	bid, err := strconv.ParseFloat(tickerResp.BidPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid price: %w", err)
	}

	ask, err := strconv.ParseFloat(tickerResp.AskPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask price: %w", err)
	}

	return &PriceData{
		Bid: bid,
		Ask: ask,
		Ts:  time.Now(),
	}, nil
}

func (c *AsterClient) FetchBTCPerpFundingRate(ctx context.Context) (*FundingRateData, error) {
	// TODO: Funding Rate 取得を実装
	return nil, fmt.Errorf("not implemented")
}
