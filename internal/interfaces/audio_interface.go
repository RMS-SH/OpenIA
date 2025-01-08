package interfaces

import (
	"context"
)

// VisionService define o contrato para operações de análise de imagens.
type AudioInterface interface {
	AudioToText(ctx context.Context, url, modelo, language string) (interface{}, error)
}
