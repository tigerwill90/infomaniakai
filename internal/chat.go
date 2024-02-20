package internal

type ChatRequest struct {
	Model             string         `json:"model,omitempty"`
	Messages          []*ChatMessage `json:"messages"`
	MaxNewTokens      int            `json:"max_new_tokens,omitempty"`
	ProfileType       string         `json:"profile_type,omitempty"`
	RepetitionPenalty float64        `json:"repetition_penalty,omitempty"`
	Seed              int            `json:"seed,omitempty"`
	SystemPrompt      string         `json:"system_prompt,omitempty"`
	Temperature       float64        `json:"temperature,omitempty"`
	TopK              int            `json:"top_k,omitempty"`
	TopP              float64        `json:"top_p,omitempty"`
	Truncate          int            `json:"truncate,omitempty"`
	TypicalP          *int           `json:"typical_p,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

type ChatResponse struct {
	Result string `json:"result,omitempty"`
	Data   struct {
		Model   string        `json:"model,omitempty"`
		Created int           `json:"created,omitempty"`
		Choices []*ChatChoice `json:"choices,omitempty"`
		Usage   struct {
			InputTokens  int `json:"input_tokens,omitempty"`
			OutputTokens int `json:"output_tokens,omitempty"`
			TotalTokens  int `json:"total_tokens,omitempty"`
		}
	} `json:"data,omitempty"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}
