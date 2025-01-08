package openia_client

import (
	"context"
	"encoding/json"
	"fmt"

	dto "github.com/RMS-SH/OpenIA/internal/dto/openia"
	"github.com/RMS-SH/OpenIA/internal/infra/clients"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// OpenAIClient implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClientText struct {
	apiKey     string
	httpClient clients.HTTPClient
	baseURL    string
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClient.
func NewOpenAIClientText(apiKey string, httpClient clients.HTTPClient) *OpenAIClientText {
	return &OpenAIClientText{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    "https://api.openai.com/v1",
	}
}

// InterpretText envia uma solicitação para interpretar o texto.
func (c *OpenAIClientText) AnalyzeText(ctx context.Context, textInput, prompt, modelo string) (interface{}, error) {
	// Monta a estrutura da requisição
	reqBody := dto.ChatCompletionsRequest{
		Model:     modelo,
		MaxTokens: 300,
		Messages: []dto.ChatMessage{
			{
				Role: "user",
				Content: []dto.MessageContent{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "text",
						Text: textInput,
					},
				},
			},
		},
	}

	// Configura headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.apiKey,
	}

	// Executa a requisição
	respBytes, err := c.httpClient.Do(ctx, "POST", fmt.Sprintf("%s/chat/completions", c.baseURL), headers, reqBody)
	if err != nil {
		return "", err
	}

	// Decodifica a resposta
	var completionResp dto.ChatCompletionsResponse
	if err := json.Unmarshal(respBytes, &completionResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da OpenAI: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("resposta da OpenAI não contém choices")
	}

	return &completionResp, nil
}

// Garante em tempo de compilação que OpenAIClient implementa VisionService.
var _ interfaces.TextInterface = (*OpenAIClientText)(nil)
