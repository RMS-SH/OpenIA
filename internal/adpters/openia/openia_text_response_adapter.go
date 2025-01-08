// internal/adapters/openia_text_response_adapter.go
package adapters_openia

import (
	"fmt"

	dto "github.com/RMS-SH/OpenIA/internal/dto/openia"
	"github.com/RMS-SH/OpenIA/internal/entities"
	errorLLM "github.com/RMS-SH/OpenIA/internal/error"
)

type OpenIAResponseAdapter struct{}

// AdaptResponse converte o objeto retornado pelo OpenIA (que supostamente
// é um *dto.ChatCompletionsResponse) em uma string única.
func (o *OpenIAResponseAdapter) AdaptResponseSupervisor(response interface{}) (string, error) {
	// Tenta converter o interface{} para *dto.ChatCompletionsResponse
	resp, ok := response.(*dto.ChatCompletionsResponse)
	if !ok {
		return "", fmt.Errorf("%w: tipo de objeto não é *dto.ChatCompletionsResponse", errorLLM.ErrAdapter)
	}

	// Verifica se existem choices
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%w: sem choices na resposta do LLM", errorLLM.ErrAdapter)
	}

	// Extrai o conteúdo do primeiro choice (ou concatena todos, etc.)

	content := resp.Choices[0].Message.Content
	return content, nil
}

// AdaptResponse converte o objeto retornado pelo OpenIA (que supostamente
// é um *dto.ChatCompletionsResponse) em uma string única.
func (o *OpenIAResponseAdapter) AdaptResponse(response interface{}) (interface{}, error) {
	// Tenta converter o interface{} para *dto.ChatCompletionsResponse
	resp, ok := response.(*dto.ChatCompletionsResponse)
	if !ok {
		return "", fmt.Errorf("%w: tipo de objeto não é *dto.ChatCompletionsResponse", errorLLM.ErrAdapter)
	}

	// Verifica se existem choices
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%w: sem choices na resposta do LLM", errorLLM.ErrAdapter)
	}

	// Extrai o conteúdo do primeiro choice (ou concatena todos, etc.)

	Resposta := entities.Text{
		Text: resp.Choices[0].Message.Content,
	}
	return Resposta, nil
}

func NewAdapterOpenIAResponseAdapter() OpenIAResponseAdapter {
	return OpenIAResponseAdapter{}
}
