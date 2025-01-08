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
	response, err := tuc.textClient.AnalyzeText(ctx, question, prompt, modelo)
	if err != nil {
		return nil, err
	}

	Adapter, err := tuc.responseAdapter.AdaptResponse(response)
	if err != nil {
		return nil, err

	}

	return Adapter, nil
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

func (uc *TextUseCase) UseCaseSupervisor(ctx context.Context, questions map[string]string, personificacaoDoModelo, modelo string) (interface{}, error) {

	type SupervisorResult struct {
		PerguntasDescartadas         []string
		PerguntasAprovadasSupervisor []string
	}

	result := SupervisorResult{
		PerguntasDescartadas:         []string{},
		PerguntasAprovadasSupervisor: []string{},
	}

	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		firstError error

		// Canal para sinalizar que devemos interromper processamento
		done = make(chan struct{})

		// sync.Once garante que algo seja executado exatamente uma vez
		once sync.Once
	)

	// Função auxiliar para fechar o canal done de forma segura
	safeCloseDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	for question, response := range questions {
		wg.Add(1)
		go func(q, r string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				// Se o contexto for cancelado, registra o primeiro erro e sai
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("contexto cancelado ou expirado: %w", ctx.Err())
				}
				mu.Unlock()
				return
			case <-done:
				// Se algum erro crítico ocorreu em outra goroutine, interrompemos aqui
				return
			default:
				// Continua o processamento normalmente
			}

			// Monta o prompt para análise de conformidade
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

			// Chamamos o cliente que faz análise de texto
			rawResponse, err := uc.textClient.AnalyzeText(ctx, prompt, personificacaoDoModelo, modelo)
			if err != nil {
				log.Printf("Erro ao analisar a pergunta '%s': %v", q, err)
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("erro ao analisar a pergunta '%s': %w", q, err)
				}
				mu.Unlock()
				safeCloseDone()
				return
			}

			// Adaptamos a resposta para extrair apenas "#SIM" ou "#NAO"
			supervisionResponse, err := uc.responseAdapter.AdaptResponseSupervisor(rawResponse)
			if err != nil {
				log.Printf("Erro ao adaptar a resposta para '%s': %v", q, err)
				mu.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("erro ao adaptar resposta da pergunta '%s': %w", q, err)
				}
				mu.Unlock()
				safeCloseDone()
				return
			}

			mu.Lock()
			defer mu.Unlock()

			// Avaliamos o conteúdo retornado
			switch supervisionResponse {
			case "#SIM":
				result.PerguntasAprovadasSupervisor = append(result.PerguntasAprovadasSupervisor, q)
			case "#NAO":
				result.PerguntasDescartadas = append(result.PerguntasDescartadas, q)
			default:
				// Se vier algo diferente de #SIM ou #NAO, consideramos erro
				log.Printf("Resposta inesperada para a pergunta '%s': %s", q, supervisionResponse)
				if firstError == nil {
					firstError = fmt.Errorf("resposta inesperada para a pergunta '%s': %s", q, supervisionResponse)
				}
				safeCloseDone()
			}
		}(question, response)
	}

	// Esperamos todas as goroutines terminarem
	wg.Wait()

	// Se ocorreu algum erro capturado, retornamos
	if firstError != nil {
		return nil, firstError
	}

	return result, nil
}
