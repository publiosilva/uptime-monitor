package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Embedder turns text into a vector for similarity search.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type OpenAIEmbedder struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewOpenAIEmbedder(apiKey, model string) *OpenAIEmbedder {
	if model == "" {
		model = "text-embedding-3-small"
	}
	return &OpenAIEmbedder{
		apiKey:     apiKey,
		model:      model,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.openai.com/v1",
	}
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if e.apiKey == "" {
		return nil, fmt.Errorf("openai api key is not configured")
	}

	body, err := json.Marshal(embeddingRequest{Model: e.model, Input: text})
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embeddings: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read embedding response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("embeddings status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed embeddingResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}
	if len(parsed.Data) == 0 {
		return nil, fmt.Errorf("embeddings returned no data")
	}
	return parsed.Data[0].Embedding, nil
}
