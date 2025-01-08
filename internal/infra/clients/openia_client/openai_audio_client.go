package openia_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"time"

	dto "github.com/RMS-SH/OpenIA/internal/dto/openia"
	"github.com/RMS-SH/OpenIA/internal/infra/clients"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
	utils "github.com/RMS-SH/UtilitariosGoRMS"
)

// OpenAIClient implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClientAudio struct {
	apiKey     string
	httpClient clients.HTTPClient
	baseURL    string
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClient.
func NewOpenAIClientAudio(apiKey string, httpClient clients.HTTPClient) *OpenAIClientAudio {
	return &OpenAIClientAudio{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    "https://api.openai.com/v1/audio/transcriptions",
	}
}

func (c *OpenAIClientAudio) AudioToText(ctx context.Context, url, modelo, language string) (interface{}, error) {
	// 1. Baixa o arquivo de áudio com timeout e validação de tamanho.
	downloadRes, err := utils.DownloadWithTimeout(url, 25, 30*time.Second)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return nil, err
	}

	audioData := downloadRes.Data
	fileName := "audio_file" + path.Ext(url)

	// 2. Prepara o multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Adiciona o arquivo de áudio
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fileWriter.Write(audioData); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	// Adiciona o campo "model"
	if err := writer.WriteField("model", modelo); err != nil {
		return nil, fmt.Errorf("failed to add model field: %w", err)
	}

	// Adiciona o campo "language"
	if err := writer.WriteField("language", language); err != nil {
		return nil, fmt.Errorf("failed to add language field: %w", err)
	}

	// Adiciona o campo "response_format"
	if err := writer.WriteField("response_format", "verbose_json"); err != nil {
		return nil, fmt.Errorf("failed to add response_format field: %w", err)
	}

	// Adiciona o campo "timestamp_granularities" sem os []
	// Use valores válidos, como "word" ou "sentence"
	if err := writer.WriteField("timestamp_granularities", "word"); err != nil {
		return nil, fmt.Errorf("failed to add timestamp_granularities field: %w", err)
	}

	// Fecha o writer para finalizar o multipart
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// 3. Prepara os headers
	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  writer.FormDataContentType(),
	}

	// 4. Cria a requisição POST com contexto
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Define os headers na requisição
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 5. Executa a requisição usando o HTTPClient
	Client := http.Client{}
	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// 6. Verifica o status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Lê a resposta de erro
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error response body: %s\n", string(errorBody)) // Log para depuração
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(errorBody))
	}

	// 7. Lê a resposta
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 8. Faz o unmarshal da resposta
	var transcriptionResponse dto.TranscriptionResponse
	if err := json.Unmarshal(respBytes, &transcriptionResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 9. Monta a estrutura de resposta com informações de download
	resposta := dto.TranscriptionResponseWithDownloadInfo{
		TranscriptionResponse: transcriptionResponse,
		DownloadIP:            downloadRes.RemoteIP,
		DownloadSizeMB:        downloadRes.SizeInMB,
		StatusCode:            resp.StatusCode,
	}

	return &resposta, nil
}

// Garante em tempo de compilação que OpenAIClient implementa VisionService.
var _ interfaces.AudioInterface = (*OpenAIClientAudio)(nil)
