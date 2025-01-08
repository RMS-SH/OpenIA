package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	adapters "github.com/RMS-SH/OpenIA/internal/adpters/openia"
	"github.com/RMS-SH/OpenIA/internal/interfaces"
)

// VisionUseCase orquestra as chamadas à VisionService.
type TextUseCase struct {
	textClient      interfaces.TextInterface
	responseAdapter adapters.OpenIAResponseAdapter
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewTextUseCase(textClient interfaces.TextInterface, responseAdapter adapters.OpenIAResponseAdapter) *TextUseCase {
	return &TextUseCase{textClient: textClient, responseAdapter: responseAdapter}
}

// AnalyzeImageFromURL encapsula a chamada de análise de imagem (URL).
func (tuc *TextUseCase) UseCasAnalyzeText(ctx context.Context, question, prompt, modelo string) (interface{}, error) {
	return tuc.textClient.AnalyzeText(ctx, question, prompt, modelo)
}

func (uc *TextUseCase) UseCasAnalyzeMultText(ctx context.Context, questions []string, prompt, modelo string) ([]interface{}, error) {

	var responses []interface{} // Corrigido o nome da variável
	for _, question := range questions {
		// Chama AnalyzeText para cada questão e coleta as respostas
		response, err := uc.textClient.AnalyzeText(ctx, question, prompt, modelo)
		if err != nil {
			return nil, err // Retorna erro caso ocorra
		}
		responses = append(responses, response) // Converte para string e adiciona ao slice
	}

	return responses, nil // Retorna todas as respostas coletadas
}

// UseCaseSupervisor processa as perguntas e respostas de forma concorrente
func (uc *TextUseCase) UseCaseSupervisor(ctx context.Context, questions map[string]string, personificacaoDoModelo, modelo string) (interface{}, error) {

	type SupervisorResult struct {
		PerguntasDescartadas         []string
		PerguntasAprovadasSupervisor []string
	}

	result := SupervisorResult{
		PerguntasDescartadas:         []string{},
		PerguntasAprovadasSupervisor: []string{},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstError error
	done := make(chan struct{})

	for question, response := range questions {
		wg.Add(1)
		go func(q, r string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("contexto cancelado ou expirado: %w", ctx.Err())
				}
				mu.Unlock()
				return
			case <-done:
				return
			default:
				// Continua o processamento
			}

			prompt := fmt.Sprintf(`

Análise de Resposta:
Você receberá dois textos:
1. **Solicitação do Usuário:** "%s"
2. **Resposta Recebida:** "%s"

Seu objetivo é verificar se a **Resposta Recebida** está em conformidade com a **Solicitação do Usuário**.
Por favor, responda apenas com:
- "#SIM" se a resposta está de acordo com o solicitado.
- "#NAO" se a resposta **não** está de acordo com o solicitado.

**Importante:** Não escreva nada além de "#SIM" ou "#NAO".
`, q, r)

			rawResponse, err := uc.textClient.AnalyzeText(ctx, prompt, personificacaoDoModelo, modelo)
			if err != nil {
				log.Printf("Erro ao analisar a pergunta '%s': %v", q, err)
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("erro ao analisar a pergunta '%s': %w", q, err)
				}
				mu.Unlock()
				close(done)
				return
			}

			supervisionResponse, err := uc.responseAdapter.AdaptResponse(rawResponse)
			if err != nil {
				log.Printf("Erro ao adaptar a resposta para '%s': %v", q, err)
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("erro ao adaptar resposta da pergunta '%s': %w", q, err)
				}
				mu.Unlock()
				close(done)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			switch supervisionResponse {
			case "#SIM":
				result.PerguntasAprovadasSupervisor = append(result.PerguntasAprovadasSupervisor, q)
			case "#NAO":
				result.PerguntasDescartadas = append(result.PerguntasDescartadas, q)
			default:
				log.Printf("Resposta inesperada para a pergunta '%s': %s", q, supervisionResponse)
				if firstError == nil {
					firstError = fmt.Errorf("resposta inesperada para a pergunta '%s': %s", q, supervisionResponse)
					close(done)
				}
			}
		}(question, response)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	return result, nil
}
