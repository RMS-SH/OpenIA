package RMSLLMs

import (
	"context"
	"errors"

	"github.com/RMS-SH/OpenIA/internal/repositories"
)

func AnalisaImagem(ctx context.Context, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem string) (string, error) {

	if llm == "OpenIA" {
		return repositories.VisionOpenIA(ctx, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem)
	}

	if llm == "Gemini" {
		return repositories.VisionGemini(ctx, llm, imageInput, apiKey, prompt, modelo, qualidadeImagem)
	}

	return "", errors.New("LLM Não informada correta!")

}
