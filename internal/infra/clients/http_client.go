package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient define uma interface para o cliente HTTP.
type HTTPClient interface {
	Do(ctx context.Context, method, url string, headers map[string]string, body interface{}) ([]byte, error)
}

// DefaultHTTPClient é a implementação padrão de HTTPClient.
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient cria e retorna uma instância de DefaultHTTPClient.
func NewDefaultHTTPClient(timeout time.Duration) *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do executa a requisição HTTP com os parâmetros fornecidos.
func (c *DefaultHTTPClient) Do(ctx context.Context, method, url string, headers map[string]string, body interface{}) ([]byte, error) {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			return nil, fmt.Errorf("erro ao codificar o corpo da requisição: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição HTTP: %w", err)
	}

	// Definir cabeçalhos
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na execução da requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("requisição HTTP falhou com status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Ler resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta HTTP: %w", err)
	}

	return respBody, nil
}
