package interfaces

import (
	"context"
)

// VisionService define o contrato para operações de análise de imagens.
type DocumentInterface interface {
	AnalyzeDocument(ctx context.Context, url, modelo string) (interface{}, error)
}
