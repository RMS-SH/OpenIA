package dto

// Este arquivo contém as estruturas para processar
// a resposta de Chat Completions da OpenAI.

// ChatCompletionsResponse é a resposta principal do endpoint /v1/chat/completions.
type ChatCompletionsResponse struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Created int64                   `json:"created"`
	Choices []chatCompletionsChoice `json:"choices"`
	// Podem existir outros campos, como "usage", etc.
}

// ChatCompletionsChoice representa cada escolha retornada pela OpenAI.
type chatCompletionsChoice struct {
	Index        int                `json:"index"`
	FinishReason string             `json:"finish_reason"`
	Message      chatCompletionsMsg `json:"message"`
}

// ChatCompletionsMsg é a mensagem retornada dentro de cada choice.
type chatCompletionsMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

///////////////////////////////////////////////////

type TranscriptionResponse struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

type transcription struct {
	Task     string                  `json:"task"`
	Language string                  `json:"language"`
	Duration float64                 `json:"duration"`
	Text     string                  `json:"text"`
	Segments []TranscriptionResponse `json:"segments"`
}
