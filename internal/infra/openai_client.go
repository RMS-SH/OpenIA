package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RMS-SH/OpenIA/internal/dto"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// OpenAIClient implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClient struct {
	apiKey string
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClient.
func NewOpenAIClient(apiKeyOpenIA string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKeyOpenIA,
	}
}

// SetAPIKey permite mudar a API Key em tempo de execução (caso seja necessário).
func (c *OpenAIClient) SetAPIKey(key string) {
	c.apiKey = key
}

// AnalyzeImage atende ao contrato da VisionService.
// Ele chama o método "AnalyzeImageInternal" que monta a requisição JSON para imagens.
func (c *OpenAIClient) AnalyzeImage(ctx context.Context, imageInput string) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("API key não foi definida")
	}

	// Exemplo: vamos usar "gpt-4o-mini" ou outro modelo que a OpenAI disponibilize para você.
	const modelName = "gpt-4o-mini"

	// Monta a estrutura base da requisição
	reqBody := dto.ChatCompletionsRequest{
		Model:     modelName,
		MaxTokens: 300,
		Messages: []dto.ChatMessage{
			{
				Role: "user",
				Content: []dto.MessageContent{
					{
						Type: "text",
						Text: "What is in this image?",
					},
				},
			},
		},
	}

	// Verifica se é URL ou Base64. Se for URL, adiciona detail:"high" (exemplo).
	if strings.HasPrefix(imageInput, "http") {
		reqBody.Messages[0].Content = append(reqBody.Messages[0].Content, dto.MessageContent{
			Type: "image_url",
			ImageURL: &dto.ImageURL{
				URL:    imageInput,
				Detail: "high",
			},
		})
	} else {
		reqBody.Messages[0].Content = append(reqBody.Messages[0].Content, dto.MessageContent{
			Type:        "image_base64",
			ImageBase64: imageInput,
		})
	}

	// Executa de fato a requisição
	return c.executeChatCompletion(ctx, reqBody)
}

// AskTextPrompt é um exemplo de outro método para o mesmo endpoint /v1/chat/completions,
// mas com mensagens de texto simples (sem "content" em array).
func (c *OpenAIClient) AskTextPrompt(ctx context.Context, prompt string) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("API key não foi definida")
	}

	// Exemplo de modelo. Ajuste conforme disponibilidade.
	const modelName = "gpt-4o"

	reqBody := dto.ChatCompletionsRequest{
		Model: modelName,
		Messages: []dto.ChatMessage{
			{
				Role:    "user",
				Content: []dto.MessageContent{{Type: "text", Text: prompt}},
			},
		},
		// max_tokens pode ser omitido ou configurado conforme necessidade
		MaxTokens: 200,
	}

	return c.executeChatCompletion(ctx, reqBody)
}

// executeChatCompletion é o método comum que faz a chamada HTTP ao endpoint da OpenAI.
// Ele pode ser reutilizado para qualquer variação de JSON que siga a estrutura base de ChatCompletionsRequest.
func (c *OpenAIClient) executeChatCompletion(ctx context.Context, request dto.ChatCompletionsRequest) (string, error) {
	// Converte para JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer marshal do corpo da requisição: %w", err)
	}

	endpoint := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Cria um HTTP Client com timeout
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na chamada à OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("chamada à OpenAI falhou com status: %d", resp.StatusCode)
	}

	// Decodifica a resposta
	var completionResp dto.ChatCompletionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da OpenAI: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("resposta da OpenAI não contém choices")
	}

	// Retorna apenas o primeiro choice
	answer := completionResp.Choices[0].Message.Content
	return answer, nil
}

// --------------------------------------------------------------------------
// Garante em tempo de compilação que OpenAIClient implementa VisionService.
// --------------------------------------------------------------------------
var _ interfaces.VisionService = (*OpenAIClient)(nil)
