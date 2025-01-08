package openia_client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RMS-SH/OpenIA/internal/dto"
	"github.com/RMS-SH/OpenIA/internal/infra/clients"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// OpenAIClient implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClientVision struct {
	apiKey     string
	httpClient clients.HTTPClient
	baseURL    string
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClient.
func NewOpenAIClientVision(apiKey string, httpClient clients.HTTPClient) *OpenAIClientVision {
	return &OpenAIClientVision{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    "https://api.openai.com/v1",
	}
}

// AnalyzeImage atende ao contrato da VisionService.
// Ele monta a requisição e utiliza o HTTPClient para fazer a chamada.
func (c *OpenAIClientVision) AnalyzeImage(ctx context.Context, imageInput, prompt, modelo, qualidadeImagem string) (interface{}, error) {
	// Monta a estrutura base da requisição
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
				},
			},
		},
	}

	// Verifica se é URL ou Base64 e adiciona o conteúdo correspondente
	if imageInput != "" {
		if strings.HasPrefix(imageInput, "http") {
			reqBody.Messages[0].Content = append(reqBody.Messages[0].Content, dto.MessageContent{
				Type: "image_url",
				ImageURL: &dto.ImageURL{
					URL:    imageInput,
					Detail: qualidadeImagem,
				},
			})
		} else {
			reqBody.Messages[0].Content = append(reqBody.Messages[0].Content, dto.MessageContent{
				Type:        "image_base64",
				ImageBase64: imageInput,
			})
		}
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
var _ interfaces.VisionService = (*OpenAIClientVision)(nil)
