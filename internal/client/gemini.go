// package client отвечает за работу с внешними API.
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

// GeminiClient — клиент для Gemini API.
type GeminiClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewGeminiClient — создание нового Gemini клиента.
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAIRecommendation — отправляет prompt в Gemini
// и возвращает текст рекомендации.
func (gc *GeminiClient) GetAIRecommendation(
	ctx context.Context,
	prompt string,
) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", errors.New("prompt is required")
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
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := gc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get recommendation")
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	candidates := result["candidates"].([]interface{})
	content := candidates[0].(map[string]interface{})
	contentMap := content["content"].(map[string]interface{})
	parts := contentMap["parts"].([]interface{})
	part := parts[0].(map[string]interface{})
	text := part["text"].(string)
	return text, nil
}