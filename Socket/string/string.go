package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SERVER_IP   = "3.88.99.255"
	SERVER_PORT = "8080"
	SERVER_ADDR = SERVER_IP + ":" + SERVER_PORT
	PROTOCOL    = "tcp"
	ALUNO_ID    = "542562"
)

func formatStringsMessage(command string, params map[string]string) string {
	args := []string{command}

	for key, value := range params {
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, "FIM")

	return strings.Join(args, "|") + "\n"
}

func sendAndReceive(conn net.Conn, message string) (string, error) {
	fmt.Printf(" -> Enviando: %s\n", strings.TrimSpace(message))

	_, err := conn.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("erro ao enviar mensagem: %w", err)
	}
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}
	return strings.TrimSpace(response), nil
}

func Authenticate(conn net.Conn, studentID string) (string, error) {
	fmt.Println("\n--- 1. AUTENTICAÇÃO (AUTH) ---")

	params := map[string]string{
		"aluno_id":  studentID,
		"timestamp": time.Now().Format("2006-01-02T15:04:05"),
	}
	authRequest := formatStringsMessage("AUTH", params)

	authResponse, err := sendAndReceive(conn, authRequest)
	if err != nil {
		return "", fmt.Errorf("falha na comunicação AUTH: %w", err)
	}

	fmt.Printf(" <- Resposta de AUTH: %s\n", authResponse)

	parts := strings.Split(authResponse, "|")
	if len(parts) < 2 || parts[0] != "OK" {
		return "", fmt.Errorf("autenticação falhou: %s", authResponse)
	}

	tokenField := parts[1]
	if !strings.HasPrefix(tokenField, "token=") {
		return "", fmt.Errorf("token não encontrado no campo de resposta: %s", tokenField)
	}

	sessionToken := strings.TrimPrefix(tokenField, "token=")
	return sessionToken, nil
}

func OperacaoEcho(conn net.Conn, token, message string) {
	fmt.Println("\n--- 2. OPERAÇÃO: ECHO ---")

	params := map[string]string{
		"operacao": "echo",
		"mensagem": message,
		"token":    token,
	}
	opRequest := formatStringsMessage("OP", params)

	opResponse, err := sendAndReceive(conn, opRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação Echo: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func OperacaoSoma(conn net.Conn, token string, numbers []int) {
	fmt.Println("\n--- 3. OPERAÇÃO: SOMA ---")

	var numsStr []string
	for _, n := range numbers {
		numsStr = append(numsStr, fmt.Sprintf("%d", n))
	}

	params := map[string]string{
		"operacao": "soma",
		"nums":     strings.Join(numsStr, ","),
		"token":    token,
	}
	opRequest := formatStringsMessage("OP", params)

	opResponse, err := sendAndReceive(conn, opRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação Soma: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func OperacaoTimestamp(conn net.Conn, token string) {
	fmt.Println("\n--- 4. OPERAÇÃO: TIMESTAMP ---")

	params := map[string]string{
		"operacao": "timestamp",
		"token":    token,
	}
	opRequest := formatStringsMessage("OP", params)

	opResponse, err := sendAndReceive(conn, opRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação Timestamp: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func OperacaoStatus(conn net.Conn, token string, detailed bool) {
	fmt.Println("\n--- 5. OPERAÇÃO: STATUS ---")

	params := map[string]string{
		"operacao": "status",
		"token":    token,
	}
	if detailed {
		params["detalhado"] = "true"
	}
	opRequest := formatStringsMessage("OP", params)

	opResponse, err := sendAndReceive(conn, opRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação Status: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func OperacaoHistorico(conn net.Conn, token string, limit int) {
	fmt.Println("\n--- 6. OPERAÇÃO: HISTÓRICO ---")

	params := map[string]string{
		"operacao": "historico",
		"token":    token,
	}
	if limit > 0 {
		params["limite"] = fmt.Sprintf("%d", limit)
	}
	opRequest := formatStringsMessage("OP", params)

	opResponse, err := sendAndReceive(conn, opRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação Histórico: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func InformacoesServidor(conn net.Conn, infoType string) {
	fmt.Println("\n--- 7. INFORMAÇÕES DO SERVIDOR (INFO) ---")

	params := map[string]string{
		"tipo": infoType,
	}
	infoRequest := formatStringsMessage("INFO", params)

	opResponse, err := sendAndReceive(conn, infoRequest)
	if err != nil {
		fmt.Printf("❌ Falha na operação INFO: %v\n", err)
		return
	}

	fmt.Printf(" <- Resposta: %s\n", opResponse)
}

func Logout(conn net.Conn, token string) {
	fmt.Println("\n--- 8. LOGOUT ---")

	params := map[string]string{
		"token": token,
	}
	logoutRequest := formatStringsMessage("LOGOUT", params)

	logoutResponse, err := sendAndReceive(conn, logoutRequest)
	if err != nil {
		fmt.Printf("⚠️ Aviso: Logout pode ter falhado (ou conexão fechada): %v\n", err)
	} else {
		fmt.Printf(" <- Resposta: %s\n", logoutResponse)
	}
}

func main() {
	fmt.Printf("--- Cliente Strings Protocol (%s) ---\n", SERVER_ADDR)

	conn, err := net.Dial(PROTOCOL, SERVER_ADDR)
	if err != nil {
		fmt.Printf("❌ Erro ao conectar a %s: %v\n", SERVER_ADDR, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("✅ Conexão TCP estabelecida.")

	sessionToken, err := Authenticate(conn, ALUNO_ID)
	if err != nil {
		fmt.Printf("❌ Erro de Autenticação: %v\n", err)
		return
	}

	fmt.Printf("✅ Autenticação bem-sucedida! Token de sessão: %s...\n", sessionToken)

	InformacoesServidor(conn, "basico")
	OperacaoEcho(conn, sessionToken, "Verificando o formato final do protocolo.")
	OperacaoSoma(conn, sessionToken, []int{10, 20, 30, 40})
	OperacaoTimestamp(conn, sessionToken)
	OperacaoStatus(conn, sessionToken, true)
	OperacaoHistorico(conn, sessionToken, 2)
	Logout(conn, sessionToken)

	fmt.Println("\n--- Cliente Strings: Fim das operações. ---")
}
