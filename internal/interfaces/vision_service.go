package interfaces

import "context"

// VisionService define o contrato para operações de análise de imagens.
type VisionService interface {
	AnalyzeImage(ctx context.Context, imageInput string) (string, error)
}
