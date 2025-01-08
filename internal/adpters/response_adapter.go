// internal/adapters/response_adapter.go
package adapters

// AnalyzeTextResponseAdapter é a interface responsável por
// converter a resposta recebida do LLM (interface{}) em uma string
// que o caso de uso entende.
type AnalyzeTextResponseAdapter interface {
	AdaptResponse(response interface{}) (string, error)
}
