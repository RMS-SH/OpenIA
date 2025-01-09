package RMSLLMs

import (
	"context"
	"errors"

	"github.com/RMS-SH/OpenIA/internal/entities"
	"github.com/RMS-SH/OpenIA/internal/repositories/openia_repositories"
)

func AnalisaImage(ctx context.Context, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem string) (*entities.ImagemAnalyzeResponse, error) {

	if llm == "OpenIA" {
		return openia_repositories.VisionOpenIA(ctx, imageInput, apiKey, prompt, modelo, qualidadeImagem)
	}
	return nil, errors.New("LLM Não informada correta!")

}

func AudioTranscription(ctx context.Context, llm, apiKey, url, modelo, language string) (*entities.AudioTranscriptionResponse, error) {

	if llm == "OpenIA" {
		return openia_repositories.AudioOpenIATranscription(ctx, apiKey, url, modelo, language)
	}
	return nil, errors.New("LLM Não informada correta!")

}

func LLMTextSimple(ctx context.Context, llm, apiKey, question, prompt, modelo string) (interface{}, error) {

	if llm == "OpenIA" {
		return openia_repositories.TextOpenIAAnalizy(ctx, question, apiKey, prompt, modelo)
	}
	return "", errors.New("LLM Não informada correta!")

}

func Supervisor(ctx context.Context, llm string, question map[string]string, apiKey, personificacaoDoModelo, modeloLLM string) (interface{}, error) {

	if llm == "OpenIA" {
		return openia_repositories.SupervisorOpenIA(ctx, question, apiKey, personificacaoDoModelo, modeloLLM)
	}
	return "", errors.New("LLM Não informada correta!")

}

func InterpretacaoPDFAssistente(ctx context.Context, prompt, url, apiKEY string) (interface{}, error) {
	return openia_repositories.InterpretacaoPDFAssistenteRepository(ctx, prompt, url, apiKEY)
}
