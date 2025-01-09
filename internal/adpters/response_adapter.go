// internal/adapters/response_adapter.go
package adapters

import "github.com/RMS-SH/OpenIA/internal/entities"

// AnalyzeTextResponseAdapter é a interface responsável por
// converter a resposta recebida do LLM (interface{}) em uma string
// que o caso de uso entende.
type AnalyzeTextResponseAdapter interface {
	AdaptResponse(response string) (interface{}, error)
	AdaptResponseSupervisor(response interface{}) (string, error)
}

type AudioToTextResponseAdapter interface {
	AdaptResponse(response interface{}) (*entities.AudioTranscriptionResponse, error)
}

type AnalyzeImageAdapter interface {
	AdaptResponse(response interface{}) (*entities.ImagemAnalyzeResponse, error)
}
