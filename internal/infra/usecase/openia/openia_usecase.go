package usecase_openia

import (
	"context"
	"fmt"
	"sync"
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

// InterpretacaoPDFAssistenteUseCase faz a interpretação de um PDF usando assistente + vector store.
// Alguns passos são executados em paralelo via go routines para melhorar a performance.
func (ac *OpenIAUseCase) InterpretacaoPDFAssistenteUseCase(ctx context.Context, prompt, url string) (interface{}, error) {
	// -------------------------------------------------------------------
	// Passo 1: Download do arquivo (fora das goroutines, pois é passo único)
	// -------------------------------------------------------------------
	downloadedPDF, err := utilitariosgorms.DownloadWithTimeout(url, 25, 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar PDF: %w", err)
	}

	// -------------------------------------------------------------------
	// Passos 2, 3 e 4 em paralelo:
	//   2. Criar Assistente
	//   3. Fazer Upload do PDF
	//   4. Criar Vector Store
	// -------------------------------------------------------------------

	var (
		wg                sync.WaitGroup
		errChan           = make(chan error, 3) // buffer para armazenar erros
		assistantResult   *openai.Assistant
		fileResult        *openai.File
		vectorStoreResult *openai.VectorStore
	)

	wg.Add(3) // Precisamos de 3 goroutines

	// Goroutine para criar assistente
	go func() {
		defer wg.Done()
		assistant, err := ac.openai.CreateAssistant(ctx, openai.GPT4oMini, "Interpretador de PDF", "Você é um assistente que interpreta o conteúdo de arquivos PDF.")
		if err != nil {
			errChan <- fmt.Errorf("erro ao criar assistente: %w", err)
			return
		}
		assistantResult = assistant
	}()

	// Goroutine para fazer upload do PDF
	go func() {
		defer wg.Done()
		file, err := ac.openai.UploadFileBytes(ctx, "Arquivo PDF.pdf", downloadedPDF.Data, openai.PurposeAssistants)
		if err != nil {
			errChan <- fmt.Errorf("erro ao fazer upload do PDF: %w", err)
			return
		}
		fileResult = file
	}()

	// Goroutine para criar vector store
	go func() {
		defer wg.Done()
		vectorStore, err := ac.openai.CreateVectorStore(ctx, "Vector Store Temporária")
		if err != nil {
			errChan <- fmt.Errorf("erro ao criar Vector Store: %w", err)
			return
		}
		vectorStoreResult = vectorStore
	}()

	// Espera todas as goroutines acima terminarem
	wg.Wait()

	// Se houve algum erro, retornamos o primeiro
	select {
	case e := <-errChan:
		return nil, e
	default:
		// Nenhum erro; segue o fluxo
	}

	// -------------------------------------------------------------------
	// Passos 5, 6 e 7 em paralelo:
	//   5. Adicionar arquivo à Vector Store
	//   6. Associar Vector Store ao Assistente (ModifyAssistant)
	//   7. Criar Thread
	// -------------------------------------------------------------------
	var (
		wg2          sync.WaitGroup
		errChan2     = make(chan error, 3)
		threadResult *openai.Thread
		// Precisamos dos resultados do passo anterior:
		//  - fileResult
		//  - assistantResult
		//  - vectorStoreResult
	)

	wg2.Add(3)

	// Goroutine 1: Adicionar arquivo na Vector Store
	go func() {
		defer wg2.Done()
		if err := ac.openai.AddFileToVectorStore(ctx, vectorStoreResult.ID, fileResult.ID); err != nil {
			errChan2 <- fmt.Errorf("erro ao adicionar arquivo à Vector Store: %w", err)
			return
		}
	}()

	// Goroutine 2: Associar Vector Store ao Assistente
	go func() {
		defer wg2.Done()
		modifiedAssistant, err := ac.openai.ModifyAssistant(ctx, assistantResult.ID, vectorStoreResult.ID)
		if err != nil {
			errChan2 <- fmt.Errorf("erro ao associar Vector Store ao Assistente: %w", err)
			return
		}
		assistantResult = modifiedAssistant // atualiza a referência
	}()

	// Goroutine 3: Criar Thread
	go func() {
		defer wg2.Done()
		thread, err := ac.openai.CreateThread(ctx, nil, nil)
		if err != nil {
			errChan2 <- fmt.Errorf("erro ao criar thread: %w", err)
			return
		}
		threadResult = thread
	}()

	wg2.Wait()

	// Se houve algum erro nesse grupo, retornamos o primeiro
	select {
	case e := <-errChan2:
		return nil, e
	default:
		// nenhum erro
	}

	// -------------------------------------------------------------------
	// Passo 8: Adicionar mensagem do usuário (depende da thread criada)
	// -------------------------------------------------------------------
	contentMsg := fmt.Sprintf(
		"%s use sua tool file_search para visualizar o pdf na vector store %s",
		prompt,
		vectorStoreResult.ID,
	)
	_, err = ac.openai.AddMessageToThread(ctx, threadResult.ID, "user", contentMsg, fileResult.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao adicionar mensagem do usuário: %w", err)
	}

	// -------------------------------------------------------------------
	// Passo 9: Criar Run
	// -------------------------------------------------------------------
	run, err := ac.openai.CreateRun(ctx, threadResult.ID, assistantResult.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a Run: %w", err)
	}

	// -------------------------------------------------------------------
	// Passo 10: Loop para Monitorar a Run
	// -------------------------------------------------------------------
	var runStatus openai.RunStatus
	for runStatus != openai.RunStatusCompleted {
		time.Sleep(5 * time.Second)
		run, err = ac.openai.RetrieveRun(ctx, threadResult.ID, run.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao recuperar status da Run: %w", err)
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
				run, err = ac.openai.SubmitToolOutputs(ctx, threadResult.ID, run.ID, toolCalls)
				if err != nil {
					return nil, fmt.Errorf("erro ao submeter tool outputs: %w", err)
				}
				fmt.Printf("Tool outputs submetidas com sucesso: %v\n", run.ID)
			}
		}
	}

	// -------------------------------------------------------------------
	// Passo 11: Recuperar mensagens
	// -------------------------------------------------------------------
	messages, err := ac.openai.ListThreadMessages(ctx, threadResult.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar mensagens: %w", err)
	}

	// -------------------------------------------------------------------
	// Passo 12: Limpeza (em paralelo), pois são independentes
	// -------------------------------------------------------------------
	var wgCleanup sync.WaitGroup
	wgCleanup.Add(4)
	cleanupErrChan := make(chan error, 4)

	// Deletar Assistente
	go func() {
		defer wgCleanup.Done()
		if err := ac.openai.DeleteAssistant(ctx, assistantResult.ID); err != nil {
			cleanupErrChan <- fmt.Errorf("erro ao deletar assistente: %w", err)
		}
	}()

	// Deletar arquivo
	go func() {
		defer wgCleanup.Done()
		if err := ac.openai.DeleteFile(ctx, fileResult.ID); err != nil {
			cleanupErrChan <- fmt.Errorf("erro ao deletar arquivo: %w", err)
		}
	}()

	// Deletar vector store
	go func() {
		defer wgCleanup.Done()
		if err := ac.openai.DeleteVectorStore(ctx, vectorStoreResult.ID); err != nil {
			cleanupErrChan <- fmt.Errorf("erro ao deletar vector store: %w", err)
		}
	}()

	// Deletar thread
	go func() {
		defer wgCleanup.Done()
		if err := ac.openai.DeleteThread(ctx, threadResult.ID); err != nil {
			cleanupErrChan <- fmt.Errorf("erro ao deletar thread: %w", err)
		}
	}()

	// Espera a limpeza
	wgCleanup.Wait()
	// Se houve algum erro na limpeza, avisamos (mas já temos as messages)
	select {
	case e := <-cleanupErrChan:
		// Decide se quer retornar erro ou apenas logar
		// Retornando, pois pode ser relevante
		return messages, e
	default:
		// sem erro
	}

	// Retornamos as mensagens resultantes da interpretação do PDF
	return messages, nil
}
