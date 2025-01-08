// internal/adapters/openia_text_response_adapter.go
package adapters

import (
	"fmt"

	"github.com/RMS-SH/OpenIA/internal/dto"
	errorLLM "github.com/RMS-SH/OpenIA/internal/error"
)

type OpenIAResponseAdapter struct{}

// AdaptResponse converte o objeto retornado pelo OpenIA (que supostamente
// é um *dto.ChatCompletionsResponse) em uma string única.
func (o *OpenIAResponseAdapter) AdaptResponse(response interface{}) (string, error) {
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

func NewAdapterOpenIAResponseAdapter() OpenIAResponseAdapter {
	return OpenIAResponseAdapter{}
}
