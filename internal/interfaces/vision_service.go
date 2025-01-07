package interfaces

import (
	"context"

	"github.com/RMS-SH/OpenIA/internal/dto"
)

// VisionService define o contrato para operações de análise de imagens.
type VisionService interface {
	AnalyzeImage(ctx context.Context, imageInput string, prompt, modelo, qualidadeImagem string) (string, error)
	AnalyzeImageFullReturn(ctx context.Context, imageInput string, prompt, modelo, qualidadeImagem string) (*dto.ChatCompletionsResponse, error)
}
