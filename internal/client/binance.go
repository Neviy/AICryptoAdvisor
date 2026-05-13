// package client отвечает за работу с внешними API.
// Здесь находится клиент для Binance API.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// BinanceClient — клиент для работы с Binance API.
type BinanceClient struct {
	baseURL string
	client  *http.Client
}

// BinancePriceResponse — ответ Binance API с ценой монеты.
type BinancePriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// NewBinanceClient — создание нового Binance клиента.
func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		baseURL: "https://api-gcp.binance.com",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetPrice — получение цены монеты по символу.
// Например: BTC -> BTCUSDT.
func (bc *BinanceClient) GetPrice(ctx context.Context,symbol string) (*BinancePriceResponse, error) {
	symbol = strings.TrimSpace(strings.ToUpper(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	url := fmt.Sprintf(
		"%s/api/v3/ticker/price?symbol=%sUSDT",
		bc.baseURL,
		symbol,
	)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := bc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"binance api returned status: %d",
			resp.StatusCode,
		)
	}
	var result BinancePriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}