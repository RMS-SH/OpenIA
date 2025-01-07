package utilits

import "github.com/RMS-SH/OpenIA/pkg/entities"

// mergeVisionOptions pega os valores do slice de overrides (opts...)
// e sobrescreve no 'base' somente os campos não-vazios.
func MergeVisionOptions(base entities.VisionOptions, overrides ...entities.VisionOptions) entities.VisionOptions {
	merged := base
	for _, o := range overrides {
		if o.Prompt != "" {
			merged.Prompt = o.Prompt
		}
		if o.Modelo != "" {
			merged.Modelo = o.Modelo
		}
		if o.QualidadeImagem != "" {
			merged.QualidadeImagem = o.QualidadeImagem
		}
	}
	return merged
}
