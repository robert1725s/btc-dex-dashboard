package dex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const hyperliquidAPIURL = "https://api.hyperliquid.xyz/info"

type HyperliquidClient struct {
	httpClient *http.Client
}

func NewHyperliquidClient() *HyperliquidClient {
	return &HyperliquidClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HyperliquidClient) Name() string {
	return "hyperliquid"
}

type hyperliquidL2Request struct {
	Type string `json:"type"`
	Coin string `json:"coin"`
}

type hyperliquidL2Response struct {
	Coin   string           `json:"coin"`
	Time   int64            `json:"time"`
	Levels [][]hyperliquidLevel `json:"levels"`
}

type hyperliquidLevel struct {
	Px string `json:"px"`
	Sz string `json:"sz"`
	N  int    `json:"n"`
}

func (c *HyperliquidClient) FetchBTCPerpPrice(ctx context.Context) (*PriceData, error) {
	reqBody := hyperliquidL2Request{
		Type: "l2Book",
		Coin: "BTC",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hyperliquidAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var l2Resp hyperliquidL2Response
	if err := json.NewDecoder(resp.Body).Decode(&l2Resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(l2Resp.Levels) < 2 || len(l2Resp.Levels[0]) == 0 || len(l2Resp.Levels[1]) == 0 {
		return nil, fmt.Errorf("invalid response: insufficient levels")
	}

	bid, err := strconv.ParseFloat(l2Resp.Levels[0][0].Px, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid price: %w", err)
	}

	ask, err := strconv.ParseFloat(l2Resp.Levels[1][0].Px, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask price: %w", err)
	}

	return &PriceData{
		Bid: bid,
		Ask: ask,
		Ts:  time.Now(),
	}, nil
}

func (c *HyperliquidClient) FetchBTCPerpFundingRate(ctx context.Context) (*FundingRateData, error) {
	// TODO: Funding Rate 取得を実装
	return nil, fmt.Errorf("not implemented")
}
