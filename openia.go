package openia

import (
	"context"
	"strings"

	"github.com/RMS-SH/OpenIA/internal/infra"
	"github.com/RMS-SH/OpenIA/internal/infra/usecase"
)

// VisionOpenIA decide se a string da imagem é URL ou Base64 e direciona para o caso de uso adequado.
func VisionOpenIA(ctx context.Context, imageInput, apiKeyOpenIA string) (string, error) {
	client := infra.NewOpenAIClient(apiKeyOpenIA)
	uc := usecase.NewVisionUseCase(client)

	// Verificação simples para distinguir URL vs. Base64
	if strings.HasPrefix(imageInput, "http") {
		return uc.AnalyzeImageFromURL(ctx, imageInput)
	}

	return uc.AnalyzeImageFromBase64(ctx, imageInput)
}
