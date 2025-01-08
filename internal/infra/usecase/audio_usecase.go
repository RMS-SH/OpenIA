package usecase

import (
	"context"

	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// VisionUseCase orquestra as chamadas à VisionService.
type AudioUseCase struct {
	audioClient interfaces.AudioInterface
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewAudioUseCase(audioClient interfaces.AudioInterface) *AudioUseCase {
	return &AudioUseCase{audioClient: audioClient}
}

// AnalyzeImageFromURL encapsula a chamada de análise de imagem (URL).
func (ac *AudioUseCase) UseCaseAudioToText(ctx context.Context, url, modelo, language string) (interface{}, error) {
	return ac.audioClient.AudioToText(ctx, url, modelo, language)
}
