package main

type BaseRequest struct {
	Tipo      string `json:"tipo"`
	AlunoID   string `json:"aluno_id,omitempty"`
	Token     string `json:"token,omitempty"`
	Timestamp string `json:"timestamp"`
}

type OperationRequest struct {
	Tipo       string `json:"tipo"`
	Operacao   string `json:"operacao"`
	Token      string `json:"token"`
	Parametros any    `json:"parametros"`
	Timestamp  string `json:"timestamp"`
}

type JsonAuthResponse struct {
	Sucesso    bool   `json:"sucesso"`
	Erro       string `json:"erro"`
	Token      string `json:"token"`
	DadosAluno struct {
		Nome      string `json:"nome"`
		Matricula string `json:"matricula"`
	} `json:"dados_aluno"`
	Mensagem  string `json:"mensagem"`
	Timestamp string `json:"timestamp"`
}

type JsonOperationResponse struct {
	Sucesso   bool           `json:"sucesso"`
	Erro      string         `json:"erro"`
	Resultado map[string]any `json:"resultado"`
	Mensagem  string         `json:"mensagem"`
	Timestamp string         `json:"timestamp"`
}

type EchoParams struct {
	Mensagem string `json:"mensagem"`
}

type SumParams struct {
	Numeros []int `json:"numeros"`
}

type TimestampParams struct{}

type StatusParams struct {
	Detalhado bool `json:"detalhado"`
}

type HistoryParams struct {
	Limite int `json:"limite"`
}
