// Exemplos de USO

package main

import (
	"fmt"

	"context"

	RMSLLMs "github.com/RMS-SH/OpenIA/internal"
)

func main() {
	// Exemplo de uso local (teste) das funções exportadas no pacote openia

	audio, _ := audioTranscription()
	fmt.Println(audio)
	imagem, _ := analisaImage()
	fmt.Println(imagem)
	text, _ := textoToLLMSimple()
	fmt.Println(text)
	supervisor, _ := supervisorRMS()
	fmt.Println(supervisor)

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
