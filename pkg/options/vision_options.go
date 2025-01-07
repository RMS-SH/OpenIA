package options

import "github.com/RMS-SH/OpenIA/pkg/entities"

// VisionOption define um tipo para opções funcionais.
type VisionOption func(*entities.VisionOptions)

// WithPrompt configura o prompt para a análise de imagem.
func WithPrompt(prompt string) VisionOption {
	return func(vo *entities.VisionOptions) {
		vo.Prompt = prompt
	}
}

// WithModelo configura o modelo a ser utilizado.
func WithModelo(modelo string) VisionOption {
	return func(vo *entities.VisionOptions) {
		vo.Modelo = modelo
	}
}

// WithQualidadeImagem configura a qualidade da imagem.
func WithQualidadeImagem(qualidade string) VisionOption {
	return func(vo *entities.VisionOptions) {
		vo.QualidadeImagem = qualidade
	}
}
