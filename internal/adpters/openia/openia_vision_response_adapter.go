// internal/adapters/openia_text_response_adapter.go
package adapters_openia

import (
	"github.com/RMS-SH/OpenIA/internal/dto"
	"github.com/RMS-SH/OpenIA/internal/entities"
)

type OpenIAResponseAdapterImagem struct{}

// AdaptResponse converte o objeto retornado pelo OpenIA (que supostamente
// é um *dto.ChatCompletionsResponse) em uma string única.
func (o *OpenIAResponseAdapterImagem) AdaptResponse(response interface{}) (interface{}, error) {
	Entrada := response.(*dto.ChatCompletionsResponse)
	Retorno := entities.ImagemAnalyzeResponse{
		Text: Entrada.Choices[0].Message.Content,
	}

	return Retorno, nil
}

func NewAdapterOpenIAResponseAdapterImagem() *OpenIAResponseAdapterImagem {
	return &OpenIAResponseAdapterImagem{}
}
