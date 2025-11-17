package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_IP   = "3.88.99.255"
	SERVER_PORT = "8081"
	SERVER_ADDR = SERVER_IP + ":" + SERVER_PORT
	PROTOCOL    = "tcp"
	ALUNO_ID    = "542562"
)

func tcpRequest(conn net.Conn, message []byte) (string, error) {
	fmt.Printf(" -> Enviando: %s\n", strings.TrimSpace(string(message)))
	_, err := conn.Write(append(message, '\n'))
	if err != nil {
		return "", fmt.Errorf("erro ao enviar: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 64*1024)

	n, err := conn.Read(buffer)

	conn.SetReadDeadline(time.Time{})

	if n > 0 {
		rawResponse := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf(" <- Recebido: %s\n", rawResponse)
		return rawResponse, nil
	}

	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	return "", fmt.Errorf("read: no data received")
}

func extractToken(authResponse string) (string, error) {
	var resp JsonAuthResponse

	if err := json.Unmarshal([]byte(authResponse), &resp); err == nil {
		if resp.Sucesso && resp.Token != "" {
			return resp.Token, nil
		}
		return "", fmt.Errorf("autenticação falhou: Token vazio")
	}
	return "", fmt.Errorf("token não encontrado ou resposta JSON inválida: %s", authResponse)
}

func Login(conn net.Conn, studentId int) (string, error) {
	fmt.Println("\n--- 1. AUTENTICAÇÃO (LOGIN) ---")
	currentTime := time.Now().Format("2006-01-02T15:04:05-07:00")
	authReq := BaseRequest{
		Tipo:      "autenticar",
		AlunoID:   strconv.Itoa(studentId),
		Timestamp: currentTime,
	}
	payload, _ := json.Marshal(authReq)
	authResponse, err := tcpRequest(conn, payload)
	if err != nil {
		return "", err
	}
	token, err := extractToken(authResponse)
	if err != nil {
		return "", fmt.Errorf("falha na autenticação: %w", err)
	}
	return token, nil
}

func RunOperation(conn net.Conn, op string, args []string, token string) {
	fmt.Printf("\n--- 2. OPERAÇÃO: %s ---\n", strings.ToUpper(op))
	var params any
	switch op {
	case "echo":
		if len(args) < 1 {
			fmt.Println("❌ echo requer mensagem")
			return
		}
		params = EchoParams{Mensagem: strings.Join(args, " ")}
	case "sum":
		if len(args) < 1 {
			fmt.Println("❌ sum requer lista")
			return
		}
		parts := strings.Split(args[0], ",")
		nums := make([]int, 0, len(parts))
		for _, p := range parts {
			if v, conv := strconv.Atoi(strings.TrimSpace(p)); conv == nil {
				nums = append(nums, v)
			}
		}
		params = SumParams{Numeros: nums}
	case "timestamp":
		params = TimestampParams{}
	case "status":
		params = StatusParams{Detalhado: true}
	case "history":
		limit := 10
		if len(args) >= 1 {
			if v, conv := strconv.Atoi(args[0]); conv == nil && v > 0 {
				limit = v
			}
		}
		params = HistoryParams{Limite: limit}
	default:
		fmt.Printf("❌ Operação desconhecida: %s\n", op)
		return
	}

	switch op {
	case "sum":
		op = "soma"
	case "history":
		op = "historico"
	}

	currentTime := time.Now().Format("2006-01-02T15:04:05-07:00")

	body := OperationRequest{
		Tipo:       "operacao",
		Operacao:   op,
		Token:      token,
		Parametros: params,
		Timestamp:  currentTime,
	}

	payload, _ := json.Marshal(body)
	rawResponse, _ := tcpRequest(conn, payload)

	var resp JsonOperationResponse
	if err := json.Unmarshal([]byte(rawResponse), &resp); err == nil {
		if resp.Sucesso {
			fmt.Printf(" ✅ Operação %s Sucesso! Resultado: %v\n", op, resp.Resultado)
		} else {
			fmt.Printf(" ⚠️ Operação %s Falha: %s\n", op, resp.Erro)
		}
	} else {
		fmt.Println("❌ Erro ao decodificar a resposta da operação.")
	}
}

func Logout(conn net.Conn, token string) {
	fmt.Println("\n--- 4. LOGOUT ---")

	currentTime := time.Now().Format("2006-01-02T15:04:05-07:00")

	logoutReq := BaseRequest{
		Tipo:      "logout",
		Token:     token,
		Timestamp: currentTime,
	}

	payload, _ := json.Marshal(logoutReq)

	_, err := tcpRequest(conn, payload)

	conn.Close()

	if err != nil {
		fmt.Printf("⚠️ Aviso: Logout pode ter falhado: %v\n", err)
	}
	fmt.Println("✅ Conexão fechada.")
}

func main() {
	fmt.Printf("--- Cliente JSON Protocol (%s) ---\n", SERVER_ADDR)

	conn, err := net.Dial(PROTOCOL, SERVER_ADDR)
	if err != nil {
		fmt.Printf("❌ Erro ao conectar a %s: %v\n", SERVER_ADDR, err)
		os.Exit(1)
	}

	fmt.Println("✅ Conexão TCP estabelecida.")

	studentIdInt, _ := strconv.Atoi(ALUNO_ID)
	sessionToken, err := Login(conn, studentIdInt)
	if err != nil {
		fmt.Printf("❌ Erro de Autenticação: %v\n", err)
		return
	}
	fmt.Printf("✅ Autenticação bem-sucedida! Token: %s...\n", sessionToken)

	RunOperation(conn, "echo", []string{"Mensagem de teste."}, sessionToken)
	RunOperation(conn, "sum", []string{"10,20,30"}, sessionToken)
	RunOperation(conn, "status", nil, sessionToken)
	RunOperation(conn, "history", []string{"5"}, sessionToken)
	RunOperation(conn, "timestamp", nil, sessionToken)

	Logout(conn, sessionToken)
	fmt.Println("\n--- Cliente JSON: Fim das operações. ---")
}
