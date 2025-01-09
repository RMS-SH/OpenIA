package usecase

import (
	"context"

	adapters "github.com/RMS-SH/OpenIA/internal/adpters"
	"github.com/RMS-SH/OpenIA/internal/entities"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// VisionUseCase orquestra as chamadas à VisionService.
type AudioUseCase struct {
	audioClient interfaces.AudioInterface
	adapter     adapters.AudioToTextResponseAdapter
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewAudioUseCase(audioClient interfaces.AudioInterface, adapter adapters.AudioToTextResponseAdapter) *AudioUseCase {
	return &AudioUseCase{audioClient: audioClient, adapter: adapter}
}

// AnalyzeImageFromURL encapsula a chamada de análise de imagem (URL).
func (ac *AudioUseCase) UseCaseAudioToText(ctx context.Context, url, modelo, language string) (*entities.AudioTranscriptionResponse, error) {

	response, err := ac.audioClient.AudioToText(ctx, url, modelo, language)
	if err != nil {
		return nil, err
	}

	Adapter, err := ac.adapter.AdaptResponse(response)
	if err != nil {
		return nil, err

	}

	return Adapter, nil
}
