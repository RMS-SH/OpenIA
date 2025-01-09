// Exemplos de USO

package main

import (
	"context"
	"fmt"

	RMSLLMs "github.com/RMS-SH/OpenIA/internal"
)

const Arquivo = "https://bot.dfktv2.com/media/whatsapp/453695414498557/1008916574347940/fozxrecife-1-pTQb.pdf"
const apiKey = "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"
const Modelo = "gpt-4o-mini"
const IDA = "file-UGvjrvrGMSsiR6WFCyFXZy"
const ID = 5
const VSID = "vs_9CWowhjE2918ca21QwZqZRpu"
const AssistantID = "asst_73pVJIHhjrCPOUNyyGBxCcJm"

func main() {
	// Exemplo de uso local (teste) das funções exportadas no pacote openia
	switch ID {
	case 1:
		audio, _ := audioTranscription()
		fmt.Println(audio)
	case 2:
		imagem, _ := analisaImage()
		fmt.Println(imagem)
	case 3:
		text, _ := textoToLLMSimple()
		fmt.Println(text)
	case 4:
		supervisor, _ := supervisorRMS()
		fmt.Println(supervisor)
	case 5:
		CreateMenssageWithFile, _ := RespostaAssistente()
		fmt.Println(CreateMenssageWithFile)

	default:
		fmt.Println("ID inválido")
	}
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

func RespostaAssistente() (interface{}, error) {
	ctx := context.Background()
	result, err := RMSLLMs.InterpretacaoPDFAssistente(ctx, "Me diga oque tem no pdf", Arquivo, apiKey)
	if err != nil {
		return nil, err

	}
	return result, nil
}
