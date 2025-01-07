package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/RMS-SH/OpenIA/internal/infra"
	"github.com/RMS-SH/OpenIA/internal/infra/usecase"
)

// LLM = OpenIA ou Gemini - Obrigatório
// imagemInPut = url ou base64 - Obrigatório
// Qualidade Imagem ( OpenIA ) = "Low", "Medium", "High"
// Função Chama Analise de Imagem com base no prompt fornecido.
func VisionGemini(ctx context.Context, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem string) (string, error) {
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
		modelo = "GPT-4o-mini"
	}
	if qualidadeImagem == "" {
		qualidadeImagem = "Low"
	}
	// Verifica se entrada da Functio Está Ok

	// Chama os Casos de USO
	client := infra.NewOpenAIClient(apiKey) // aqui é um exemplo porém mudaria apenas para GeminiCliente
	uc := usecase.NewVisionUseCase(client)

	// Verificação simples para distinguir URL vs. Base64
	if strings.HasPrefix(imageInput, "http") {
		return uc.AnalyzeImageFromURL(ctx, imageInput, prompt, modelo, qualidadeImagem)
	}

	return uc.AnalyzeImageFromBase64(ctx, imageInput, prompt, modelo, qualidadeImagem)
}
