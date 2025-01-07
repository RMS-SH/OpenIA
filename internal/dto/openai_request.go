package dto

// Este arquivo contém as estruturas para montar o JSON
// de requisições ao endpoint de Chat Completions.

// ChatCompletionsRequest modela a requisição geral de chat completions para OpenAI.
type ChatCompletionsRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
	// Podem existir outros campos: "tools", "functions", "temperature", etc.
	// dependendo do tipo de requisição e parâmetros desejados.
}

// ChatMessage modela cada mensagem enviada no chat.
type ChatMessage struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content,omitempty"`
	// Para casos de mensagens somente com texto simples (sem array de content),
	// poderíamos ter um campo "content" string. A API do Chat Completions
	// aceita diferentes formatos.
}

// MessageContent define o conteúdo que pode ser texto ou imagem.
// Aqui adicionamos "detail" para o caso do "image_url" que pode ter "detail":"high".
type MessageContent struct {
	Type        string    `json:"type"`
	Text        string    `json:"text,omitempty"`
	ImageURL    *ImageURL `json:"image_url,omitempty"`
	ImageBase64 string    `json:"image_base64,omitempty"`
}

// ImageURL modela a estrutura interna para URLs de imagem,
// podendo conter parâmetros adicionais, como "detail".
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}
