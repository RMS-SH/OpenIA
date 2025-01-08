package openia_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	dto "github.com/RMS-SH/OpenIA/internal/dto/openia"
	utilitariosgorms "github.com/RMS-SH/UtilitariosGoRMS"
)

// OpenAIClient implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClientAssistent struct {
	apiKey  string
	baseURL string
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClient.
func NewOpenAIClientAssistent(apiKey string) *OpenAIClientAssistent {
	return &OpenAIClientAssistent{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1/assistants",
	}
}

func (ass *OpenAIClientAssistent) CadastraAssistenteSimples(ctx context.Context, model, prompt string) (interface{}, error) {
	// Serializa assistantRequest para JSON

	RequestModel := dto.CriaAssistenteSimplesDTO(model, prompt)

	reqBody, err := json.Marshal(RequestModel)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %w", err)
	}

	// Monta a requisição HTTP

	req, err := http.NewRequest(http.MethodPost, ass.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request HTTP: %w", err)
	}

	// Adiciona os headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ass.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Cria um client HTTP com timeout
	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	// Executa a requisição
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao executar request HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta HTTP: %w", err)
	}

	// Verifica se não foi 2xx
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("status não esperado: %d. Resposta: %s", resp.StatusCode, buf.String())
	}

	// Retorna o JSON de resposta como string
	return buf.String(), nil
}

func (ass *OpenAIClientAssistent) DeletaAssistent(ctx context.Context, assistantID string) (interface{}, error) {
	// Monta a URL completa com o ID do assistente
	url := fmt.Sprintf("%s/%s", ass.baseURL, assistantID)

	// Cria a requisição HTTP DELETE
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	// Adiciona os headers necessários
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ass.apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Cria um client HTTP com timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Executa a requisição
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica o status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Lê a resposta para obter detalhes do erro
		buf := new(bytes.Buffer)
		_, readErr := buf.ReadFrom(resp.Body)
		if readErr != nil {
			return nil, err
		}
		return nil, err
	}

	// Se necessário, pode processar a resposta aqui.
	// Geralmente, DELETE retorna 204 No Content, então pode não haver corpo.
	return "Ok", nil
}

func (c *OpenAIClientAssistent) UploadFile(ctx context.Context, fileURL string) (interface{}, error) {
	// Faz o download do arquivo usando a função DownloadWithTimeout
	Arquivo, err := utilitariosgorms.DownloadWithTimeout(fileURL, 25, 30*time.Second) // supondo que retorne DownloadedFile
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar o arquivo: %w", err)
	}

	// Cria um buffer para armazenar o corpo da requisição
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Adiciona o campo "purpose"
	err = writer.WriteField("purpose", "assistants") // ou "assistants" conforme necessário
	if err != nil {
		return nil, fmt.Errorf("erro ao adicionar campo 'purpose': %w", err)
	}

	// Cria o campo do arquivo com o nome do arquivo
	part, err := writer.CreateFormFile("file", "temp")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar campo de arquivo: %w", err)
	}

	// Copia o conteúdo do arquivo para o campo
	_, err = io.Copy(part, bytes.NewReader(Arquivo.Data)) // Arquivo.Data é []byte
	if err != nil {
		return nil, fmt.Errorf("erro ao copiar conteúdo do arquivo: %w", err)
	}

	// Fecha o writer para finalizar a escrita do multipart
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("erro ao fechar writer multipart: %w", err)
	}

	// Monta a URL completa diretamente
	uploadURL := fmt.Sprintf("%s/files", "https://api.openai.com/v1")

	// Cria a requisição HTTP POST
	req, err := http.NewRequest("POST", uploadURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona os headers necessários
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	Client := http.Client{}
	// Executa a requisição
	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta HTTP: %w", err)
	}

	// Verifica o status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status não esperado: %d. Resposta: %s", resp.StatusCode, string(respBody))
	}

	// Deserializa a resposta
	var uploadResp dto.UploadFileResponse
	err = json.Unmarshal(respBody, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar resposta: %w", err)
	}

	return &uploadResp, nil
}

// Estruturas de Requisição e Resposta
type CreateVectorStoreRequest struct {
	FileIDs      []string             `json:"file_ids,omitempty"`
	Name         string               `json:"name,omitempty"`
	ExpiresAfter *ExpiresAfterRequest `json:"expires_after,omitempty"`
	// Outros campos opcionais podem ser adicionados aqui
}

type ExpiresAfterRequest struct {
	Anchor string `json:"anchor"` // Campo obrigatório
	Days   int    `json:"days"`   // Campo obrigatório
}

type CreateVectorStoreResponse struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	CreatedAt  int64  `json:"created_at"`
	Name       string `json:"name"`
	Bytes      int64  `json:"bytes"`
	FileCounts struct {
		InProgress int `json:"in_progress"`
		Completed  int `json:"completed"`
		Failed     int `json:"failed"`
		Cancelled  int `json:"cancelled"`
		Total      int `json:"total"`
	} `json:"file_counts"`
	// Outros campos da resposta podem ser adicionados aqui
}

type APIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (c *OpenAIClientAssistent) CreateVectorStore(name, id string) (interface{}, error) {

	// Validação dos campos obrigatórios
	if name == "" {
		return nil, fmt.Errorf("nome da vector store é obrigatório")
	}
	if id == "" {
		return nil, fmt.Errorf("file ID é obrigatório")
	}

	// Cria o corpo da requisição com o file_id fornecido
	reqBody := CreateVectorStoreRequest{
		FileIDs: []string{id},
		Name:    name,
		ExpiresAfter: &ExpiresAfterRequest{
			Anchor: "last_active_at", // Campo obrigatório
			Days:   1,                // Campo obrigatório
		},
	}

	// Serializa o corpo para JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar o corpo da requisição: %v", err)
	}

	// Cria a requisição HTTP
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/vector_stores", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	// Define os headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	// Inicializa o cliente HTTP
	client := &http.Client{}

	// Envia a requisição
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Tenta deserializar a resposta de erro
		var apiErr APIErrorResponse
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("requisição falhou com status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("requisição falhou com status %d: %s", resp.StatusCode, apiErr.Error.Message)
	}

	// Deserializa a resposta JSON
	var response CreateVectorStoreResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar a resposta: %v", err)
	}

	return &response, nil
}
