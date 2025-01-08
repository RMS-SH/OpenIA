package errorLLM

import "fmt"

// ErroAdapter indica que a conversão do objeto não foi bem sucedida.
var ErrAdapter = fmt.Errorf("falha ao adaptar objeto de resposta")
