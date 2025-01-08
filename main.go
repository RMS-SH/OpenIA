// Exemplos de USO

package main

import (
	"fmt"

	"context"

	RMSLLMs "github.com/RMS-SH/OpenIA/internal"
)

const arquivo = "https://bot.dfktv2.com/media/whatsapp/453695414498557/1258988938489402/2024-modelo-de-contrato-de-compra-e-venda-de-passagens-aereas-HXn0.docx"
const apiKey = "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"
const Modelo = "gpt-4o-mini"
const IDA = "file-UGvjrvrGMSsiR6WFCyFXZy"

func main() {
	// Exemplo de uso local (teste) das funções exportadas no pacote openia

	// audio, _ := audioTranscription()
	// fmt.Println(audio)
	// imagem, _ := analisaImage()
	// fmt.Println(imagem)
	// text, _ := textoToLLMSimple()
	// fmt.Println(text)
	// supervisor, _ := supervisorRMS()
	// fmt.Println(supervisor)
	// cadastraAssistente, _ := cadatraAssistente()
	// fmt.Println(cadastraAssistente)
	// deletaAssistente, _ := deletaAssistente()
	// fmt.Println(deletaAssistente)
	// uploadArquivo, _ := uploadArquivoOpenIA()
	// fmt.Println(uploadArquivo)
	vs, _ := vs()
	fmt.Println(vs)

}

func audioTranscription() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.AudioTranscription(
		ctx,
		"OpenIA",
		"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",
		"https://bot.dfktv2.com/media/whatsapp/453695414498557/1258612962086712/-cTWi.mp3",
		"",
		"",
	)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func analisaImage() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.AnalisaImage(
		ctx,
		"OpenIA",
		"https://bot.dfktv2.com/media/whatsapp/453695414498557/1001385551828554/-5Nr2.png",
		"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",
		"Analise bem direito a foto",
		"gpt-4o-mini",
		"",
	)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func textoToLLMSimple() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.LLMTextSimple(
		ctx,
		"OpenIA",
		"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",
		"Olá tudo bem?",
		"atue como atendente de uma usina nuclear",
		"gpt-4o-mini",
	)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func supervisorRMS() (interface{}, error) {
	vars := make(map[string]string, 1)
	vars["Quero um bolo"] = "Certo tome um pastel"
	ctx := context.Background()
	result, err := RMSLLMs.Supervisor(
		ctx,
		"OpenIA",
		vars,
		"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",
		"Atue como atendente especilizado em padaria",
		"",
	)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func cadatraAssistente() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.CadastraAssistenteOpenIA(ctx, apiKey, Modelo, "Olá mundo")
	if err != nil {
		return nil, err

	}
	return result, nil
}

func deletaAssistente() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.DeletaAssistentOpenIA(ctx, apiKey, "asst_uftbVSjcj9svm7X94IzHgdiZ")
	if err != nil {
		return nil, err

	}
	return result, nil
}

func uploadArquivoOpenIA() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.UpdaloadArquivoOpenaIA(ctx, arquivo, apiKey)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func vs() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.CreateVectorStoreByFileID(ctx, IDA, apiKey)
	if err != nil {
		return nil, err

	}
	return result, nil
}
