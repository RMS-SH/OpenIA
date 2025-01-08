package RMSLLMs

import (
	"context"
	"errors"

	"github.com/RMS-SH/OpenIA/internal/repositories/openia_repositories"
)

func AnalisaImage(ctx context.Context, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem string) (interface{}, error) {

	if llm == "OpenIA" {
		return openia_repositories.VisionOpenIA(ctx, imageInput, apiKey, prompt, modelo, qualidadeImagem)
	}
	return "", errors.New("LLM Não informada correta!")

}

func AudioTranscription(ctx context.Context, llm, apiKey, url, modelo, language string) (interface{}, error) {

	if llm == "OpenIA" {
		return openia_repositories.AudioOpenIATranscription(ctx, apiKey, url, modelo, language)
	}
	return "", errors.New("LLM Não informada correta!")

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

func CadastraAssistenteOpenIA(ctx context.Context, apiKey, modelo, prompt string) (interface{}, error) {
	return openia_repositories.CadastraAssisnteSimples(ctx, apiKey, modelo, prompt)
}

func DeletaAssistentOpenIA(ctx context.Context, apiKey, id string) (interface{}, error) {
	return openia_repositories.DeletaAssistentSimples(ctx, id, apiKey)
}

func UpdaloadArquivoOpenaIA(ctx context.Context, url, apiKEY string) (interface{}, error) {
	return openia_repositories.UpdaloadArquivoOpenIA(ctx, url, apiKEY)
}

func CreateVectorStoreByFileID(ctx context.Context, id, apiKEY string) (interface{}, error) {
	return openia_repositories.CreateVectorStoreByArquive(ctx, id, apiKEY)
}
