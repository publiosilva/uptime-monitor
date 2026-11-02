package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SQLGenerator produces a read-only SQL query for a natural-language question.
// Implementations must return an empty string when the question cannot be answered
// from users/monitors/heartbeats.
type SQLGenerator interface {
	GenerateSQL(ctx context.Context, question string) (string, error)
}

const sqlGeneratorSystemPrompt = `You write a single PostgreSQL SELECT (or WITH ... SELECT) for an uptime monitoring database.
Tables:
- users(id UUID, email TEXT, password_hash TEXT, created_at TIMESTAMPTZ)
- monitors(id UUID, user_id UUID, name VARCHAR, url TEXT, method VARCHAR, timeout INT, frequency INT, is_active BOOLEAN, is_up BOOLEAN, created_at TIMESTAMPTZ)
- heartbeats(id BIGSERIAL, monitor_id UUID, status_code INT, latency_ms INT, is_up BOOLEAN, error_message TEXT, created_at TIMESTAMPTZ)

Rules:
- Output ONLY the SQL, no markdown fences or commentary.
- The query MUST filter to the authenticated user with parameter $1 (the user UUID). Example: monitors.user_id = $1 OR join monitors and filter user_id = $1.
- Read-only: SELECT or WITH ... SELECT only. Never INSERT/UPDATE/DELETE/DDL.
- If the question cannot be answered from these tables (e.g. weather), or asks for another user's data you cannot scope with $1 alone, return exactly: UNANSWERABLE
`

type OpenAISQLGenerator struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewOpenAISQLGenerator(apiKey, model string) *OpenAISQLGenerator {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAISQLGenerator{
		apiKey:     apiKey,
		model:      model,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.openai.com/v1",
	}
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (g *OpenAISQLGenerator) GenerateSQL(ctx context.Context, question string) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("openai api key is not configured")
	}

	body, err := json.Marshal(chatCompletionRequest{
		Model: g.model,
		Messages: []chatMessage{
			{Role: "system", Content: sqlGeneratorSystemPrompt},
			{Role: "user", Content: question},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("chat completions: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read chat response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("chat completions status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("chat completions returned no choices")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	content = stripSQLFences(content)
	if strings.EqualFold(content, "UNANSWERABLE") {
		return "", nil
	}
	return content, nil
}

func stripSQLFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```sql")
		s = strings.TrimPrefix(s, "```SQL")
		s = strings.TrimPrefix(s, "```")
		if idx := strings.LastIndex(s, "```"); idx >= 0 {
			s = s[:idx]
		}
	}
	return strings.TrimSpace(s)
}
