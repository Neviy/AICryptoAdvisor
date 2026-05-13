// package client- отвечает за работу с внешними API.
// Здесь находится клиент для Gemini API.
package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type GeminiClient struct {
	apiKey		string
	baseURL	string
	client *http.Client
}

type AIRecommendation struct {
	Text string `json:"text"`
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,	
		baseURL: "https://generativelanguage.googleapis.com",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (gc *GeminiClient) GetAIRecommendation(ctx context.Context,prompt string) (*AIRecommendation, error) {
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}
	url := gc.baseURL +
		"/v1beta/models/gemini-2.0-flash:generateContent?key=" +
		gc.apiKey
	body := `{
		"contents": [
			{
				"parts": [
					{
						"text": "` + prompt + `"
					}
				]
			}
		]
	}`

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		strings.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get recommendation")
	}
	// decode response
	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	candidates := result["candidates"].([]interface{})
	content := candidates[0].(map[string]interface{})
	contentMap := content["content"].(map[string]interface{})
	parts := contentMap["parts"].([]interface{})
	part := parts[0].(map[string]interface{})
	text := part["text"].(string)
	return &AIRecommendation{
		Text: text,
	}, nil
}