package main

import (
	"fmt"
	"log"

	"context"

	RMSLLMs "github.com/RMS-SH/OpenIA/internal"
)

func main() {
	// Exemplo de uso local (teste) das funções exportadas no pacote openia
	ctx := context.Background()

	result, err := RMSLLMs.AnalisaImagemRetornoCompleto(ctx, "OpenIA", "https://bot.dfktv2.com/media/whatsapp/453695414498557/1001385551828554/-5Nr2.png", "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", "Analise bem direito a foto", "gpt-4o-mini", "")
	if err != nil {
		log.Fatal("Erro ao analisar imagem:", err)
	}
	fmt.Println("[VisionOpenIA] Resultado:", result)

	// ...
	fmt.Println("Aplicação executada com sucesso!")
}
