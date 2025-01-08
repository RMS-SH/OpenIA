package dto

/*
Estrutura principal para enviar ao endpoint /v1/assistants

Considere que alguns campos são ponteiros para permitir omissão (omitempty) no JSON.
Em Go, valores padrão (ex.: string = "") ainda são serializados,
então usar ponteiros nos dá controle para enviar somente quando necessário.
*/
type CreateAssistantRequest struct {
	Model          string            `json:"model"`                     // Obrigatório
	Name           *string           `json:"name,omitempty"`            // Opcional
	Description    *string           `json:"description,omitempty"`     // Opcional
	Instructions   *string           `json:"instructions,omitempty"`    // Opcional
	Tools          []Tool            `json:"tools,omitempty"`           // Opcional
	Metadata       map[string]string `json:"metadata,omitempty"`        // Opcional
	Temperature    *float64          `json:"temperature,omitempty"`     // Opcional
	TopP           *float64          `json:"top_p,omitempty"`           // Opcional
	ResponseFormat interface{}       `json:"response_format,omitempty"` // Opcional (pode ser "auto" ou objeto complexo)
}

/*
Estrutura genérica de Tool.
Dependendo do "type", as propriedades específicas ficam em sub-estruturas.
*/
type Tool struct {
	Type       string        `json:"type"`                  // "code_interpreter", "file_search" ou "function"
	FileSearch *FileSearch   `json:"file_search,omitempty"` // Somente para type="file_search"
	Function   *FunctionTool `json:"function,omitempty"`    // Somente para type="function"
	// Em caso de "code_interpreter", não há campos adicionais além de "type"
}

/*
Estrutura para file_search tool, com configurações adicionais se necessário.
*/
type FileSearch struct {
	// Exemplo de possíveis overrides específicos de file_search
	// Ajuste de acordo com a documentação das 'overrides' para file_search.
	// Aqui está apenas ilustrativo:
	MaxResults *int `json:"max_results,omitempty"`
}

/*
Estrutura para a ferramenta type="function".
*/
type FunctionTool struct {
	Name        string                 `json:"name"`                  // Obrigatório
	Description *string                `json:"description,omitempty"` // Opcional
	Parameters  map[string]interface{} `json:"parameters,omitempty"`  // JSON Schema parcial
	Strict      *bool                  `json:"strict,omitempty"`      // Se quiser usar schema estrito
}

/*
Estrutura para o response_format do tipo json_schema (exemplo).
Pode ser expandida conforme necessidade.
*/
type JSONSchemaResponseFormat struct {
	Type       string     `json:"type"`        // "json_schema"
	JSONSchema JSONSchema `json:"json_schema"` // Objeto que define a estrutura
}

type JSONSchema struct {
	Name   string                 `json:"name"`             // Nome do schema
	Schema map[string]interface{} `json:"schema,omitempty"` // JSON schema
	Strict *bool                  `json:"strict,omitempty"` // Se deseja validar estritamente
	Descr  *string                `json:"description,omitempty"`
}

/*
Funções auxiliares para criar Tools específicos
*/

// Cria um tool do tipo code_interpreter (sem campos adicionais).
func NewCodeInterpreterTool() Tool {
	return Tool{
		Type: "code_interpreter",
	}
}

// Cria um tool do tipo file_search com possibilidade de overrides.
func NewFileSearchTool(maxResults int) Tool {
	return Tool{
		Type: "file_search",
		FileSearch: &FileSearch{
			MaxResults: &maxResults,
		},
	}
}

// Cria um tool do tipo function, com parâmetros personalizáveis.
func NewFunctionTool(name string, description string, strict bool, parameters map[string]interface{}) Tool {
	return Tool{
		Type: "function",
		Function: &FunctionTool{
			Name:        name,
			Description: &description,
			Parameters:  parameters,
			Strict:      &strict,
		},
	}
}

/*
Função para criar o corpo de requisição (CreateAssistantRequest).
Você pode customizar apenas os campos que desejar.
*/
func NewCreateAssistantRequest(
	model string,
	name *string,
	description *string,
	instructions *string,
	tools []Tool,
	metadata map[string]string,
	temperature *float64,
	topP *float64,
	responseFormat interface{},
) CreateAssistantRequest {

	return CreateAssistantRequest{
		Model:          model,
		Name:           name,
		Description:    description,
		Instructions:   instructions,
		Tools:          tools,
		Metadata:       metadata,
		Temperature:    temperature,
		TopP:           topP,
		ResponseFormat: responseFormat,
	}
}

// Exemplo de Cadastro de Function Tools

// functionTool := NewFunctionTool(
// 	"calcular_soma",
// 	"Função que retorna a soma de dois números inteiros",
// 	true,
// 	map[string]interface{}{
// 		"type": "object",
// 		"properties": map[string]interface{}{
// 			"a": map[string]interface{}{
// 				"type": "integer",
// 			},
// 			"b": map[string]interface{}{
// 				"type": "integer",
// 			},
// 		},
// 		"required": []string{"a", "b"},
// 	},
// )

// boolPtr é apenas um helper para criar ponteiro de bool
func boolPtr(b bool) *bool {
	return &b
}

func CriaAssistenteSimplesDTO(model, instructions string) *CreateAssistantRequest {

	// Exemplo de metadata
	metadata := map[string]string{
		"projeto":   "RMS File Interpreter",
		"autor":     "Oswaldo",
		"categoria": "IA",
	}

	// Exemplo de temperature e top_p
	temp := 1.0
	topP := 1.0

	// Montamos a requisição
	assistantReq := NewCreateAssistantRequest(
		model,
		nil,
		nil,
		&instructions,
		nil,
		metadata,
		&temp,
		&topP,
		"auto", // coloque "auto" se quiser o padrão
	)

	return &assistantReq
}

type UploadFileResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// Estrutura de requisição para criação de vector store
type CreateVectorStoreRequest struct {
	FileIDs          []string               `json:"file_ids,omitempty"`
	Name             string                 `json:"name,omitempty"`
	ExpiresAfter     map[string]interface{} `json:"expires_after,omitempty"`
	ChunkingStrategy map[string]interface{} `json:"chunking_strategy,omitempty"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

// Estrutura de resposta para criação de vector store
type CreateVectorStoreResponse struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	CreatedAt  int64  `json:"created_at"`
	Name       string `json:"name"`
	Bytes      int    `json:"bytes"`
	FileCounts struct {
		InProgress int `json:"in_progress"`
		Completed  int `json:"completed"`
		Failed     int `json:"failed"`
		Cancelled  int `json:"cancelled"`
		Total      int `json:"total"`
	} `json:"file_counts"`
}
