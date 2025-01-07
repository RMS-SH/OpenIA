package usecase

import (
	"context"

	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// VisionUseCase orquestra as chamadas à VisionService.
type VisionUseCase struct {
	visionClient interfaces.VisionService
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewVisionUseCase(visionClient interfaces.VisionService) *VisionUseCase {
	return &VisionUseCase{visionClient: visionClient}
}

// AnalyzeImageFromURL encapsula a chamada de análise de imagem (URL).
func (uc *VisionUseCase) AnalyzeImageFromURL(ctx context.Context, url string) (string, error) {
	return uc.visionClient.AnalyzeImage(ctx, url)
}

// AnalyzeImageFromBase64 encapsula a chamada de análise de imagem (Base64).
func (uc *VisionUseCase) AnalyzeImageFromBase64(ctx context.Context, base64 string) (string, error) {
	return uc.visionClient.AnalyzeImage(ctx, base64)
}
