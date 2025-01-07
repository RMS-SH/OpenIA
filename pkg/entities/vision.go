package entities

// VisionOptions define as opções para análise de imagem.
type VisionOptions struct {
	Prompt          string
	Modelo          string
	QualidadeImagem string
}

// VisionRequest encapsula os dados necessários para processar uma análise de imagem.
type VisionRequest struct {
	ImageInput   string
	ApiKeyOpenIA string
	Options      VisionOptions
}
