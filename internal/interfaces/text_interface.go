package interfaces

import (
	"context"
)

// VisionService define o contrato para operações de análise de imagens.
type TextInterface interface {
	AnalyzeText(ctx context.Context, text string, prompt, modelo string) (interface{}, error)
}
