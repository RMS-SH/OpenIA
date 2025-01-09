package openia_client

import (
	"context"
	"fmt"
	"time"

	utilitariosgorms "github.com/RMS-SH/UtilitariosGoRMS"
	openai "github.com/sashabaranov/go-openai"
)

// OpenAIClientAssistent implementa a interface VisionService e lida com chamadas à API da OpenAI.
type OpenAIClientAssistent struct {
	ClientOpenAI *openai.Client
}

// NewOpenAIClient cria e retorna uma instância de OpenAIClientAssistent.
func NewOpenAIClientAssistent(ClientOpenAI *openai.Client) *OpenAIClientAssistent {
	return &OpenAIClientAssistent{
		ClientOpenAI: ClientOpenAI,
	}
}

/*
==============================
   FUNÇÕES PARA ASSISTENTES
==============================
*/

// CreateAssistant cria um novo assistente
func (ass *OpenAIClientAssistent) CreateAssistant(ctx context.Context, model, name, instructions string) (*openai.Assistant, error) {
	request := openai.AssistantRequest{
		Model:        model,
		Name:         &name,
		Instructions: &instructions,
		Tools: []openai.AssistantTool{
			{
				Type: openai.AssistantToolTypeFileSearch,
			},
		},
	}

	resp, err := ass.ClientOpenAI.CreateAssistant(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar assistente: %w", err)
	}
	return &resp, nil
}

// ModifyAssistant atualiza um assistente existente (ex: adicionar VectorStoreID)
func (ass *OpenAIClientAssistent) ModifyAssistant(ctx context.Context, assistantID string, request openai.AssistantRequest) (*openai.Assistant, error) {
	resp, err := ass.ClientOpenAI.ModifyAssistant(ctx, assistantID, request)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar assistente: %w", err)
	}
	return &resp, nil
}

// DeleteAssistant deleta um assistente
func (ass *OpenAIClientAssistent) DeleteAssistant(ctx context.Context, assistantID string) error {
	_, err := ass.ClientOpenAI.DeleteAssistant(ctx, assistantID)
	if err != nil {
		return fmt.Errorf("erro ao deletar assistente: %w", err)
	}
	return nil
}

/*
==============================
     FUNÇÕES PARA ARQUIVOS
==============================
*/

// UploadFileBytes faz upload de um arquivo (em bytes)
func (ass *OpenAIClientAssistent) UploadFileBytes(ctx context.Context, fileName string, content []byte, purpose openai.PurposeType) (*openai.File, error) {
	fileReq := openai.FileBytesRequest{
		Name:    fileName,
		Purpose: purpose,
		Bytes:   content,
	}
	resp, err := ass.ClientOpenAI.CreateFileBytes(ctx, fileReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do arquivo: %w", err)
	}
	return &resp, nil
}

// DownloadAndUploadFile auxilia no download do arquivo de uma URL e upload em seguida
func (ass *OpenAIClientAssistent) DownloadAndUploadFile(ctx context.Context, fileURL string, timeoutSeconds int, purpose openai.PurposeType) (*openai.File, error) {
	// Exemplo de uso da lib interna para download
	byteArquivo, err := utilitariosgorms.DownloadWithTimeout(fileURL, 25, time.Duration(timeoutSeconds)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar o arquivo: %w", err)
	}

	return ass.UploadFileBytes(ctx, fileURL, byteArquivo.Data, purpose)
}

// DeleteFile exclui um arquivo da OpenAI
func (ass *OpenAIClientAssistent) DeleteFile(ctx context.Context, fileID string) error {
	err := ass.ClientOpenAI.DeleteFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("erro ao deletar arquivo: %w", err)
	}
	return nil
}

/*
==============================
  FUNÇÕES PARA VECTOR STORE
==============================
*/

// CreateVectorStore cria uma Vector Store vazia
func (ass *OpenAIClientAssistent) CreateVectorStore(ctx context.Context, name string) (*openai.VectorStore, error) {
	vectorStoreRequest := openai.VectorStoreRequest{
		Name: name,
	}
	resp, err := ass.ClientOpenAI.CreateVectorStore(ctx, vectorStoreRequest)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar Vector Store: %w", err)
	}
	return &resp, nil
}

// AddFileToVectorStore associa um arquivo a uma Vector Store existente
func (ass *OpenAIClientAssistent) AddFileToVectorStore(ctx context.Context, vectorStoreID, fileID string) error {
	req := openai.VectorStoreFileRequest{
		FileID: fileID,
	}
	_, err := ass.ClientOpenAI.CreateVectorStoreFile(ctx, vectorStoreID, req)
	if err != nil {
		return fmt.Errorf("erro ao adicionar arquivo na Vector Store: %w", err)
	}
	return nil
}

// DeleteVectorStore exclui uma Vector Store
func (ass *OpenAIClientAssistent) DeleteVectorStore(ctx context.Context, vectorStoreID string) error {
	_, err := ass.ClientOpenAI.DeleteVectorStore(ctx, vectorStoreID)
	if err != nil {
		return fmt.Errorf("erro ao deletar o vector store: %w", err)
	}
	return nil
}

/*
==============================
    FUNÇÕES PARA THREAD
==============================
*/

// CreateThread cria uma nova thread (sem mensagens iniciais ou podendo inserir alguma de boas-vindas, se desejar)
func (ass *OpenAIClientAssistent) CreateThread(ctx context.Context, messages []openai.ThreadMessage, metadata map[string]any) (*openai.Thread, error) {
	req := openai.ThreadRequest{
		Messages: messages,
		Metadata: metadata,
	}
	thread, err := ass.ClientOpenAI.CreateThread(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a thread: %w", err)
	}
	return &thread, nil
}

// DeleteThread deleta uma thread
func (ass *OpenAIClientAssistent) DeleteThread(ctx context.Context, threadID string) error {
	_, err := ass.ClientOpenAI.DeleteThread(ctx, threadID)
	if err != nil {
		return fmt.Errorf("erro ao deletar a thread: %w", err)
	}
	return nil
}

// AddMessageToThread adiciona uma mensagem a uma thread existente
func (ass *OpenAIClientAssistent) AddMessageToThread(ctx context.Context, threadID, role, content, fileID string) (*openai.Message, error) {
	msgReq := openai.MessageRequest{
		Role:    role,
		Content: content,
	}

	// Se tiver um fileID, adicionamos como attachment
	if fileID != "" {
		msgReq.Attachments = []openai.ThreadAttachment{
			{
				FileID: fileID,
				Tools: []openai.ThreadAttachmentTool{
					{Type: "file_search"},
				},
			},
		}
	}

	resp, err := ass.ClientOpenAI.CreateMessage(ctx, threadID, msgReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a mensagem: %w", err)
	}
	return &resp, nil
}

// ListThreadMessages recupera as mensagens de uma thread
func (ass *OpenAIClientAssistent) ListThreadMessages(ctx context.Context, threadID string) (*openai.MessagesList, error) {
	messages, err := ass.ClientOpenAI.ListMessage(ctx, threadID, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar mensagens: %w", err)
	}
	return &messages, nil
}

/*
==============================
        FUNÇÕES DE RUN
==============================
*/

// CreateRun cria uma nova Run para o assistente dentro de uma thread
func (ass *OpenAIClientAssistent) CreateRun(ctx context.Context, threadID, assistantID string) (*openai.Run, error) {
	runReq := openai.RunRequest{
		AssistantID: assistantID,
	}
	run, err := ass.ClientOpenAI.CreateRun(ctx, threadID, runReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a run: %w", err)
	}
	return &run, nil
}

// RetrieveRun recupera o status de uma Run
func (ass *OpenAIClientAssistent) RetrieveRun(ctx context.Context, threadID, runID string) (*openai.Run, error) {
	run, err := ass.ClientOpenAI.RetrieveRun(ctx, threadID, runID)
	if err != nil {
		return nil, fmt.Errorf("erro ao recuperar a run: %w", err)
	}
	return &run, nil
}

// SubmitToolOutputs envia as saídas das tools (funções) chamadas pelo assistente
func (ass *OpenAIClientAssistent) SubmitToolOutputs(ctx context.Context, threadID, runID string, toolCalls []openai.ToolCall) (*openai.Run, error) {
	var toolOutputs []openai.ToolOutput
	for _, tc := range toolCalls {
		if tc.Type == "function" {
			// Geralmente, para file_search, pode ser só "{}" de retorno
			toolOutputs = append(toolOutputs, openai.ToolOutput{
				ToolCallID: tc.ID,
				Output:     "{}",
			})
		} else {
			// Se for necessário tratar de outra forma
			resp, _ := handleToolCall(tc)
			toolOutputs = append(toolOutputs, openai.ToolOutput{
				ToolCallID: tc.ID,
				Output:     resp,
			})
		}
	}

	run, err := ass.ClientOpenAI.SubmitToolOutputs(ctx, threadID, runID, openai.SubmitToolOutputsRequest{
		ToolOutputs: toolOutputs,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao submeter tool outputs: %w", err)
	}
	return &run, nil
}

/*
==============================

	EXEMPLO DE USO DA FUNÇÃO
	handleToolCall

==============================
*/
func handleToolCall(toolCall openai.ToolCall) (string, error) {
	fmt.Printf("Function invoked: %s\n, arguments: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
	// Implemente a lógica para lidar com a chamada de função aqui, se necessário.
	return "{}", nil // Retorne a resposta da função como uma string JSON.
}
