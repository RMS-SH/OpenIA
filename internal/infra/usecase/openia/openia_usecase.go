package usecase_openia

import (
	"context"
	"fmt"
	"time"

	interfaces_openia "github.com/RMS-SH/OpenIA/internal/interfaces/openia"
	utilitariosgorms "github.com/RMS-SH/UtilitariosGoRMS"
	openai "github.com/sashabaranov/go-openai"
)

// VisionUseCase orquestra as chamadas à VisionService.
type OpenIAUseCase struct {
	openai interfaces_openia.OpenIAInterface
}

// NewVisionUseCase injeta a interface que contém a implementação real (OpenAIClient).
func NewOpenIAUseCase(openai interfaces_openia.OpenIAInterface) *OpenIAUseCase {
	return &OpenIAUseCase{openai: openai}
}

func (ac *OpenIAUseCase) InterpretacaoPDFAssistenteUseCase(ctx context.Context, prompt, url string) (interface{}, error) {

	DownloadPDF, err := utilitariosgorms.DownloadWithTimeout(url, 25, 25*time.Second)
	if err != nil {
		return nil, err
	}

	// 1. Criar Assistente
	assistant, err := ac.openai.CreateAssistant(ctx, openai.GPT4oMini, "Interpretador de PDF", "Você é um assistente que interpreta o conteúdo de arquivos PDF.")
	if err != nil {
		return nil, err
	}

	// 2. Fazer Upload do PDF
	file, err := ac.openai.UploadFileBytes(ctx, "Arquivo PDF.pdf", DownloadPDF.Data, openai.PurposeAssistants)
	if err != nil {
		return nil, err
	}

	// 3. Criar Vector Store
	vectorStore, err := ac.openai.CreateVectorStore(ctx, "Vector Store Temporária")
	if err != nil {
		return nil, err
	}

	// 4. Adicionar arquivo à Vector Store
	if err := ac.openai.AddFileToVectorStore(ctx, vectorStore.ID, file.ID); err != nil {
		return nil, err
	}

	// 5. Associar Vector Store ao Assistente
	modifyReq := openai.AssistantRequest{
		Model:   "gpt-4o-mini",
		FileIDs: []string{vectorStore.ID}, // ou vectorStoreID dependendo da API
	}
	assistant, err = ac.openai.ModifyAssistant(ctx, assistant.ID, modifyReq)
	if err != nil {
		return nil, err
	}

	// 6. Criar Thread
	thread, err := ac.openai.CreateThread(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	// 7. Adicionar mensagem do usuário
	contentMsg := fmt.Sprintf("%s use sua tool file_search para visualizar o pdf na vector store %s", prompt, vectorStore.ID)
	_, err = ac.openai.AddMessageToThread(ctx, thread.ID, "user", contentMsg, file.ID)
	if err != nil {
		return nil, err
	}

	// 8. Criar Run
	run, err := ac.openai.CreateRun(ctx, thread.ID, assistant.ID)
	if err != nil {
		return nil, err
	}

	// 9. Loop para Monitorar Run
	var runStatus openai.RunStatus
	for runStatus != openai.RunStatusCompleted {
		time.Sleep(5 * time.Second)
		run, err = ac.openai.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			return nil, err
		}
		runStatus = run.Status
		fmt.Printf("Status da Run: %s\n", runStatus)

		switch runStatus {
		case openai.RunStatusFailed, openai.RunStatusCancelled, openai.RunStatusExpired:
			return nil, fmt.Errorf("run falhou: %v", run.LastError)
		case openai.RunStatusRequiresAction:
			if run.RequiredAction.Type == openai.RequiredActionTypeSubmitToolOutputs {
				fmt.Println("Handling tool calls ...")
				toolCalls := run.RequiredAction.SubmitToolOutputs.ToolCalls
				run, err = ac.openai.SubmitToolOutputs(ctx, thread.ID, run.ID, toolCalls)
				if err != nil {
					return nil, err
				}
				fmt.Printf("Tool outputs submetidas com sucesso: %v\n", run.ID)
			}
		}
	}

	// 10. Recuperar mensagens
	messages, err := ac.openai.ListThreadMessages(ctx, thread.ID)
	if err != nil {
		return nil, err
	}

	// 11. Limpeza: Deletar assistente, arquivo, vector store e thread
	if err := ac.openai.DeleteAssistant(ctx, assistant.ID); err != nil {
		return nil, err
	}
	if err := ac.openai.DeleteFile(ctx, file.ID); err != nil {
		return nil, err
	}
	if err := ac.openai.DeleteVectorStore(ctx, vectorStore.ID); err != nil {
		return nil, err
	}
	if err := ac.openai.DeleteThread(ctx, thread.ID); err != nil {
		return nil, err
	}

	return messages, nil

}
