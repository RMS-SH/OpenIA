package usecase_openia

import (
	"context"

	interfaces_openia "github.com/RMS-SH/OpenIA/internal/interfaces/openia"
)

// VisionUseCase orquestra as chamadas à VisionService.
type OpenIAUseCase struct {
	openia interfaces_openia.OpenIAInterface
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewOpenIAUseCase(openia interfaces_openia.OpenIAInterface) *OpenIAUseCase {
	return &OpenIAUseCase{openia: openia}
}

func (ac *OpenIAUseCase) OpenIACreateAssistent(ctx context.Context, modelo, language string) (interface{}, error) {

	response, err := ac.openia.CadastraAssistenteSimples(ctx, modelo, language)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ac *OpenIAUseCase) ExcluirAssistent(ctx context.Context, id string) (interface{}, error) {

	response, err := ac.openia.DeletaAssistent(ctx, id)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ac *OpenIAUseCase) UpdaloadArquivo(ctx context.Context, url string) (interface{}, error) {

	response, err := ac.openia.UploadFile(ctx, url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ac *OpenIAUseCase) CreateVectorStore(ctx context.Context, id string) (interface{}, error) {

	response, err := ac.openia.CreateVectorStore("TEMP", id)
	if err != nil {
		return nil, err
	}

	return response, nil
}
