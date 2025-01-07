package dto

// Este arquivo contém as estruturas para processar
// a resposta de Chat Completions da OpenAI.

// ChatCompletionsResponse é a resposta principal do endpoint /v1/chat/completions.
type ChatCompletionsResponse struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Created int64                   `json:"created"`
	Choices []ChatCompletionsChoice `json:"choices"`
	// Podem existir outros campos, como "usage", etc.
}

// ChatCompletionsChoice representa cada escolha retornada pela OpenAI.
type ChatCompletionsChoice struct {
	Index        int                `json:"index"`
	FinishReason string             `json:"finish_reason"`
	Message      ChatCompletionsMsg `json:"message"`
}

// ChatCompletionsMsg é a mensagem retornada dentro de cada choice.
type ChatCompletionsMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
