// internal/adapters/openia_text_response_adapter.go
package adapters_openia

import (
	dto "github.com/RMS-SH/OpenIA/internal/dto/openia"
	"github.com/RMS-SH/OpenIA/internal/entities"
)

type OpenIAResponseAdapterAudio struct{}

// AdaptResponse converte o objeto retornado pelo OpenIA (que supostamente
// é um *dto.ChatCompletionsResponse) em uma string única.
func (o *OpenIAResponseAdapterAudio) AdaptResponse(response interface{}) (interface{}, error) {
	Entrada := response.(*dto.TranscriptionResponseWithDownloadInfo)
	Retorno := entities.AudioTranscriptionResponse{
		Text:             Entrada.TranscriptionResponse.Text,
		DurationSegundos: Entrada.TranscriptionResponse.Duration,
	}

	return Retorno, nil
}

func NewAdapterOpenIAResponseAdapterAudio() *OpenIAResponseAdapterAudio {
	return &OpenIAResponseAdapterAudio{}
}
