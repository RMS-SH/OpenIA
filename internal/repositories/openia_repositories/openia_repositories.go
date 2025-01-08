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
	usecase_openia "github.com/RMS-SH/OpenIA/internal/infra/usecase/openia"
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
	adapter := adapters.NewAdapterOpenIAResponseAdapterImagem()
	uc := usecase.NewVisionUseCase(client, adapter)

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
		language = "pt"
	}
	if modelo == "" {
		modelo = "whisper-1"
	}

	httpClient := clients.NewDefaultHTTPClient(360 * time.Second)

	// Chama os Casos de USO
	client := openia_client.NewOpenAIClientAudio(apiKey, httpClient)
	adapter := adapters.NewAdapterOpenIAResponseAdapterAudio()
	uc := usecase.NewAudioUseCase(client, adapter)

	return uc.UseCaseAudioToText(ctx, url, modelo, language)
}

func SupervisorOpenIA(ctx context.Context, question map[string]string, apiKey, personificacaoDoModelo, modeloLLM string) (interface{}, error) {
	// Verificação obrigatória dos parâmetros
	if question == nil {
		return "", errors.New("question não pode ser vazio")
	}

	if apiKey == "" {
		return "", errors.New("API KEY não pode ser vazia")
	}

	if personificacaoDoModelo == "" {
		personificacaoDoModelo = "Atue como supervisor"
	}
	if modeloLLM == "" {
		modeloLLM = "gpt-4o-mini"
	}

	httpClient := clients.NewDefaultHTTPClient(360 * time.Second)

	// Chama os Casos de USO
	client := openia_client.NewOpenAIClientText(apiKey, httpClient)
	adapter := adapters.NewAdapterOpenIAResponseAdapter()
	uc := usecase.NewTextUseCase(client, adapter)

	return uc.UseCaseSupervisor(ctx, question, personificacaoDoModelo, modeloLLM)
}

func CadastraAssisnteSimples(ctx context.Context, apiKey, modelo, prompt string) (interface{}, error) {

	client := openia_client.NewOpenAIClientAssistent(apiKey)
	uc := usecase_openia.NewOpenIAUseCase(client)

	return uc.OpenIACreateAssistent(ctx, modelo, prompt)

}

func DeletaAssistentSimples(ctx context.Context, id, apiKey string) (interface{}, error) {

	client := openia_client.NewOpenAIClientAssistent(apiKey)
	uc := usecase_openia.NewOpenIAUseCase(client)

	return uc.ExcluirAssistent(ctx, id)

}

func UpdaloadArquivoOpenIA(ctx context.Context, url, apiKEY string) (interface{}, error) {

	client := openia_client.NewOpenAIClientAssistent(apiKEY)
	uc := usecase_openia.NewOpenIAUseCase(client)

	return uc.UpdaloadArquivo(ctx, url)

}

func CreateVectorStoreByArquive(ctx context.Context, id, apiKEY string) (interface{}, error) {

	client := openia_client.NewOpenAIClientAssistent(apiKEY)
	uc := usecase_openia.NewOpenIAUseCase(client)

	return uc.CreateVectorStore(ctx, id)

}
