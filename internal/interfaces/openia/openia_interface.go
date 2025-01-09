package interfaces_openia

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

// OpenIAInterface define o contrato para operações com a API da OpenAI
type OpenIAInterface interface {
	// Funções para Assistentes
	CreateAssistant(ctx context.Context, model, name, instructions string) (*openai.Assistant, error)
	ModifyAssistant(ctx context.Context, assistantID string, vectorStoreID string) (*openai.Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID string) error

	// Funções para Arquivos
	UploadFileBytes(ctx context.Context, fileName string, content []byte, purpose openai.PurposeType) (*openai.File, error)
	DownloadAndUploadFile(ctx context.Context, fileURL string, timeoutSeconds int, purpose openai.PurposeType) (*openai.File, error)
	DeleteFile(ctx context.Context, fileID string) error

	// Funções para Vector Store
	CreateVectorStore(ctx context.Context, name string) (*openai.VectorStore, error)
	AddFileToVectorStore(ctx context.Context, vectorStoreID, fileID string) error
	DeleteVectorStore(ctx context.Context, vectorStoreID string) error

	// Funções para Thread
	CreateThread(ctx context.Context, messages []openai.ThreadMessage, metadata map[string]any) (*openai.Thread, error)
	DeleteThread(ctx context.Context, threadID string) error
	AddMessageToThread(ctx context.Context, threadID, role, content, fileID string) (*openai.Message, error)
	ListThreadMessages(ctx context.Context, threadID string) (*openai.MessagesList, error)

	// Funções de Run
	CreateRun(ctx context.Context, threadID, assistantID string) (*openai.Run, error)
	RetrieveRun(ctx context.Context, threadID, runID string) (*openai.Run, error)
	SubmitToolOutputs(ctx context.Context, threadID, runID string, toolCalls []openai.ToolCall) (*openai.Run, error)
}
