package infomaniakai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tigerwill90/infomaniakai/internal"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"

	"github.com/tmc/langchaingo/llms"
)

const defaultBaseUrl = "https://api.infomaniak.com"

type ChatMessage = internal.ChatMessage

type LLM struct {
	CallbacksHandler callbacks.Handler
	url              string
	c                *http.Client
}

const (
	RoleAssistant = "assistant"
	RoleUser      = "user"
)

var _ llms.Model = (*LLM)(nil)

// Call implements llms.Model.
func (o *LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, o, prompt, options...)
}

// GenerateContent implements llms.Model.
func (o *LLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMGenerateContentStart(ctx, messages)
	}

	opts := llms.CallOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	chatMsgs := make([]*ChatMessage, 0, len(messages))
	systemPrompt := ""
	for _, mc := range messages {
		msg := &ChatMessage{Content: mc.Parts[0].(llms.TextContent).Text}

		switch mc.Role {
		case schema.ChatMessageTypeSystem:
			systemPrompt = msg.Content
			continue
		case schema.ChatMessageTypeAI:
			msg.Role = RoleAssistant
		case schema.ChatMessageTypeHuman:
			msg.Role = RoleUser
		case schema.ChatMessageTypeGeneric:
			msg.Role = RoleUser
		default:
			return nil, fmt.Errorf("role %s not supported", mc.Role)
		}

		chatMsgs = append(chatMsgs, msg)
	}

	chatReq := &internal.ChatRequest{
		Model:             opts.Model,
		Messages:          chatMsgs,
		MaxNewTokens:      opts.MaxTokens,
		ProfileType:       "standard",
		RepetitionPenalty: opts.RepetitionPenalty,
		Seed:              opts.Seed,
		SystemPrompt:      systemPrompt,
		Temperature:       opts.Temperature,
		TopK:              opts.TopK,
		TopP:              opts.TopP,
		Truncate:          5000,
		TypicalP:          nil,
	}

	buf, err := json.Marshal(chatReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chatResp internal.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	response := &llms.ContentResponse{
		Choices: make([]*llms.ContentChoice, 0, len(chatResp.Data.Choices)),
	}

	for _, choice := range chatResp.Data.Choices {
		response.Choices = append(response.Choices, &llms.ContentChoice{
			Content:    choice.Message.Content,
			StopReason: choice.FinishReason,
		})
	}

	return response, nil
}

func New(opts ...Option) (*LLM, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &LLM{
		c: &http.Client{
			Transport: NewAuthTransport(cfg.key, http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
		url: defaultBaseUrl + fmt.Sprintf("/1/llm/%d", cfg.productId),
	}, nil
}

type AuthTransport struct {
	rt  http.RoundTripper
	key string
}

func (t AuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", "Bearer "+t.key)
	return t.rt.RoundTrip(request)
}

func NewAuthTransport(key string, rt http.RoundTripper) http.RoundTripper {
	return &AuthTransport{
		rt:  rt,
		key: key,
	}
}
