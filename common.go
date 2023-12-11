package mistral

type UsageInfo struct {
	PromptTokens     int  `json:"prompt_tokens"`
	TotalTokens      int  `json:"total_tokens"`
	CompletionTokens *int `json:"completion_tokens"`
}
