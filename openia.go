// openia.go
package openia

import (
	"context"
	"errors"
	"strings"

	"github.com/RMS-SH/OpenIA/internal/infra"
	"github.com/RMS-SH/OpenIA/internal/infra/usecase"
	"github.com/RMS-SH/OpenIA/pkg/entities"
	"github.com/RMS-SH/OpenIA/pkg/utilits"
)

func VisionOpenIA(ctx context.Context, imageInput, apiKeyOpenIA string, opts ...entities.VisionOptions) (string, error) {
	// Verificação obrigatória dos parâmetros
	if imageInput == "" {
		return "", errors.New("imageInput não pode ser vazio")
	}

	if apiKeyOpenIA == "" {
		return "", errors.New("API KEY não pode ser vazia")
	}

	// Configurar valores padrão
	options := entities.VisionOptions{
		Prompt:          "Analise a Imagem",
		Modelo:          "GPT-4o-mini",
		QualidadeImagem: "low",
	}

	// Mescla os valores de cada VisionOptions do slice `opts` no 'options' base
	options = utilits.MergeVisionOptions(options, opts...)

	client := infra.NewOpenAIClient(apiKeyOpenIA)
	uc := usecase.NewVisionUseCase(client)

	// Verificação simples para distinguir URL vs. Base64
	if strings.HasPrefix(imageInput, "http") {
		return uc.AnalyzeImageFromURL(ctx, imageInput, options.Prompt, options.Modelo, options.QualidadeImagem)
	}

	return uc.AnalyzeImageFromBase64(ctx, imageInput, options.Prompt, options.Modelo, options.QualidadeImagem)
}
