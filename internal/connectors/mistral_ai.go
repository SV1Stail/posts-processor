package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	Endpoint          = "https://api.mistral.ai"
	CodestralEndpoint = "https://codestral.mistral.ai"
	DefaultMaxRetries = 5
	DefaultTimeout    = 120 * time.Second
	mistralTiny       = "mistral-tiny"

	roleUser      = "user"
	roleAssistant = "assistant"
	roleSystem    = "system"

	FinishReasonStop        FinishReason = "stop"
	FinishReasonLength      FinishReason = "length"
	FinishReasonModelLength FinishReason = "model_length"
)

var retryStatusCodes = map[int]bool{
	429: true,
	500: true,
	502: true,
	503: true,
	504: true,
}

type MistralClient struct {
	apiKey     string
	endpoint   string
	maxRetries int
	timeout    time.Duration
	model      string
}

type NewMistralClientRequest struct {
	ApiKey     string
	Endpoint   string
	MaxRetries int
	Timeout    time.Duration
	Model      string
}

func NewMistralClient(in *NewMistralClientRequest) *MistralClient {
	if in.ApiKey == "" {
		in.ApiKey = os.Getenv("MISTRAL_API_KEY")
	}
	if in.Endpoint == "" {
		in.Endpoint = Endpoint
	}
	if in.MaxRetries == 0 {
		in.MaxRetries = DefaultMaxRetries
	}
	if in.Timeout == 0 {
		in.Timeout = DefaultTimeout
	}
	if in.Model == "" {
		in.Model = mistralTiny
	}

	return &MistralClient{
		apiKey:     in.ApiKey,
		endpoint:   in.Endpoint,
		maxRetries: in.MaxRetries,
		timeout:    in.Timeout,
		model:      in.Model,
	}
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
}

// MistralChatResponse описывает полный ответ от API чата Mistral.
type MistralChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice представляет один из вариантов ответа модели.
type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`
	FinishReason FinishReason `json:"finish_reason"`
}

// Message — роль и содержимое сообщения.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage содержит информацию о количестве токенов.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FinishReason — причина остановки генерации.
type FinishReason string

func (mc *MistralClient) NewChat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	url := mc.endpoint + "/v1/chat/completions"
	reqBody := ChatRequest{
		Model: mc.model,
		Messages: []Message{
			{
				Role:    roleSystem,
				Content: systemPrompt,
			},
			{
				Role:    roleUser,
				Content: userMessage,
			},
		},
		Temperature: 0.5,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("marshal failed")
		return "", err
	}

	client := &http.Client{Timeout: mc.timeout}

	var resp *http.Response
	for i := 0; i < mc.maxRetries; i++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Err(err).Ctx(ctx).Msg("create request failed")

			return "", err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+mc.apiKey)

		resp, err = client.Do(httpReq)
		if err != nil {
			log.Warn().Ctx(ctx).
				Int("attempt", i).
				Msg("request failed")
			if i == mc.maxRetries-1 {
				log.Err(err).Ctx(ctx).
					Msg("request failed")
				return "", err
			}
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				continue
			}
		}
		if _, ok := retryStatusCodes[resp.StatusCode]; ok {
			resp.Body.Close()
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			}

			continue
		}

		break
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("read body failed")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("mistral API returned status %d: %s", resp.StatusCode, string(body))
		log.Err(err).Ctx(ctx).Msg("status code not 200")

		return "", err
	}

	var chatResp MistralChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		log.Err(err).Ctx(ctx).Msg("unmarshal resp failed")
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		err = fmt.Errorf("no choices in response")
		log.Err(err).Ctx(ctx).Msg("status code not 200")

		return "", err
	}

	log.Info().Ctx(ctx).Interface("response", chatResp).Msg("Mistral resp")

	return chatResp.Choices[0].Message.Content, nil
}
