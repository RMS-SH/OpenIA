package openia_repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	adapters "github.com/RMS-SH/OpenIA/internal/adpters/openia"
	"github.com/RMS-SH/OpenIA/internal/infra/clients"
	"github.com/RMS-SH/OpenIA/internal/infra/clients/openia_client"
	"github.com/RMS-SH/OpenIA/internal/infra/usecase"
)

// imagemInPut = url ou base64 - Obrigatório
// Qualidade Imagem ( OpenIA ) = "Low", "Medium", "High"
// Função Chama Analise de Imagem com base no prompt fornecido.
func VisionOpenIA(ctx context.Context, imageInput, apiKey, prompt, modelo, qualidadeImagem string) (interface{}, error) {
	// Verificação obrigatória dos parâmetros
	if imageInput == "" {
		return "", errors.New("imageInput não pode ser vazio")
	}

	if apiKey == "" {
		return "", errors.New("API KEY não pode ser vazia")
	}

	if prompt == "" {
		prompt = "Analise a Imagem"
	}
	if modelo == "" {
		modelo = "gpt-4o-mini"
	}
	if qualidadeImagem == "" {
		qualidadeImagem = "low"
	}

	httpClient := clients.NewDefaultHTTPClient(360 * time.Second)

	// Chama os Casos de USO
	client := openia_client.NewOpenAIClientVision(apiKey, httpClient)
	uc := usecase.NewVisionUseCase(client)

	// Verificação simples para distinguir URL vs. Base64
	if strings.HasPrefix(imageInput, "http") {
		return uc.UseCaseAnalyzeImageFromURL(ctx, imageInput, prompt, modelo, qualidadeImagem)
	}

	return uc.UseCasAnalyzeImageFromBase64(ctx, imageInput, prompt, modelo, qualidadeImagem)
}

func TextOpenIAAnalizy(ctx context.Context, question, apiKey, prompt, modelo string) (interface{}, error) {
	// Verificação obrigatória dos parâmetros
	if question == "" {
		return "", errors.New("question não pode ser vazio")
	}

	if apiKey == "" {
		return "", errors.New("API KEY não pode ser vazia")
	}

	if prompt == "" {
		prompt = "Analise a Imagem"
	}
	if modelo == "" {
		modelo = "gpt-4o-mini"
	}

	httpClient := clients.NewDefaultHTTPClient(360 * time.Second)

	// Chama os Casos de USO
	client := openia_client.NewOpenAIClientText(apiKey, httpClient)
	adapter := adapters.NewAdapterOpenIAResponseAdapter()
	uc := usecase.NewTextUseCase(client, adapter)

	return uc.UseCasAnalyzeText(ctx, question, prompt, modelo)
}

func AudioOpenIATranscription(ctx context.Context, apiKey, url, modelo, language string) (interface{}, error) {
	// Verificação obrigatória dos parâmetros
	if url == "" {
		return "", errors.New("question não pode ser vazio")
	}

	if apiKey == "" {
		return "", errors.New("API KEY não pode ser vazia")
	}

	if language == "" {
		language = "Portuguese"
	}
	if modelo == "" {
		modelo = "whisper-1"
	}

	httpClient := clients.NewDefaultHTTPClient(360 * time.Second)

	// Chama os Casos de USO
	client := openia_client.NewOpenAIClientAudio(apiKey, httpClient)
	uc := usecase.NewAudioUseCase(client)

	return uc.UseCaseAudioToText(ctx, url, "", "")
}