package openia_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path"

	"github.com/RMS-SH/OpenIA/internal/dto"
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

// InterpretText envia uma solicitação para interpretar o texto.
func (c *OpenAIClientAudio) AudioToText(ctx context.Context, url, modelo, language string) (interface{}, error) {

	err := utils.FileSizeFromURLVerify(url, 25)
	if err != nil {
		return nil, err
	}

	audioData, err := utils.DownloadFileFromURLToBinary(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return nil, err
	}

	fileName := "audio_file" + path.Ext(url)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fileWriter.Write(audioData); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	// Add additional parameters
	if err := writer.WriteField("language", language); err != nil {
		return nil, fmt.Errorf("failed to add language field: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"Authorization":   "Bearer " + c.apiKey,
		"Content-Type":    writer.FormDataContentType(),
		"response_format": "verbose_json",
	}

	// Execute request
	respBytes, err := c.httpClient.Do(ctx, "POST", fmt.Sprintf("%s", c.baseURL), headers, body)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	var transcriptionResponse dto.TranscriptionResponse
	if err := json.Unmarshal(respBytes, &transcriptionResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &transcriptionResponse, nil

}

// Garante em tempo de compilação que OpenAIClient implementa VisionService.
var _ interfaces.AudioInterface = (*OpenAIClientAudio)(nil)
