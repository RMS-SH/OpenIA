package interfaces_openia

import (
	"context"
)

// VisionService define o contrato para operações de análise de imagens.
type OpenIAInterface interface {
	CadastraAssistenteSimples(ctx context.Context, model, prompt string) (interface{}, error)
	DeletaAssistent(ctx context.Context, assistantID string) (interface{}, error)
	UploadFile(ctx context.Context, fileURL string) (interface{}, error)
	CreateVectorStore(name, id string) (interface{}, error)
}
