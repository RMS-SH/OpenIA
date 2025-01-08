package usecase

import (
	"context"

	adapters "github.com/RMS-SH/OpenIA/internal/adpters"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// VisionUseCase orquestra as chamadas à VisionService.
type VisionUseCase struct {
	visionClient  interfaces.VisionService
	adapterVision adapters.AnalyzeImageAdapter
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewVisionUseCase(visionClient interfaces.VisionService, adapterVision adapters.AnalyzeImageAdapter) *VisionUseCase {
	return &VisionUseCase{visionClient: visionClient, adapterVision: adapterVision}
}

// AnalyzeImageFromURL encapsula a chamada de análise de imagem (URL).
func (uc *VisionUseCase) UseCaseAnalyzeImageFromURL(ctx context.Context, url, prompt, modelo, qualidadeImagem string) (interface{}, error) {

	response, err := uc.visionClient.AnalyzeImage(ctx, url, prompt, modelo, qualidadeImagem)
	if err != nil {
		return nil, err
	}

	Adapter, err := uc.adapterVision.AdaptResponse(response)
	if err != nil {
		return nil, err

	}

	return Adapter, nil

}

// AnalyzeImageFromBase64 encapsula a chamada de análise de imagem (Base64).
func (uc *VisionUseCase) UseCasAnalyzeImageFromBase64(ctx context.Context, base64, prompt, modelo, qualidadeImagem string) (interface{}, error) {
	response, err := uc.visionClient.AnalyzeImage(ctx, base64, prompt, modelo, qualidadeImagem)
	if err != nil {
		return nil, err
	}
	Adapter, err := uc.adapterVision.AdaptResponse(response)
	if err != nil {
		return nil, err

	}

	return Adapter, nil
}
