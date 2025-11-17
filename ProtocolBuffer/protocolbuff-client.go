package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	pb "pb/proto"

	"google.golang.org/protobuf/proto"
)

const (
	SERVER_IP   = "3.88.99.255"
	SERVER_PORT = "8082"
	SERVER_ADDR = SERVER_IP + ":" + SERVER_PORT
	PROTOCOL    = "tcp"
	ALUNO_ID    = "542562"
)

func encodeMessage(m proto.Message) ([]byte, error) {
	payload, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar Protobuf: %w", err)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payload)))

	return append(header, payload...), nil
}

func decodeResponse(conn net.Conn) (*pb.Resposta, error) {
	reader := bufio.NewReader(conn)

	header := make([]byte, 4)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("erro ao ler o cabeçalho (4 bytes): %w", err)
	}

	payloadSize := binary.BigEndian.Uint32(header)
	if payloadSize == 0 {
		return nil, fmt.Errorf("tamanho do payload inválido (0)")
	}

	payload := make([]byte, payloadSize)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, fmt.Errorf("erro ao ler o payload (%d bytes): %w", payloadSize, err)
	}

	resp := &pb.Resposta{}
	if err := proto.Unmarshal(payload, resp); err != nil {
		return nil, fmt.Errorf("erro ao deserializar Protobuf: %w", err)
	}

	return resp, nil
}

func Authenticate(conn net.Conn, alunoID string) (string, error) {
	fmt.Println("\n--- 1. AUTENTICAÇÃO (AUTH) ---")

	authReq := &pb.Requisicao{
		Tipo: &pb.Requisicao_Auth{
			Auth: &pb.Auth{
				AlunoId: alunoID,
			},
		},
	}

	encodedMsg, err := encodeMessage(authReq)
	if err != nil {
		return "", fmt.Errorf("falha ao encodar AUTH: %w", err)
	}

	if _, err := conn.Write(encodedMsg); err != nil {
		return "", fmt.Errorf("falha ao enviar AUTH: %w", err)
	}
	fmt.Printf(" -> Enviando %d bytes...\n", len(encodedMsg))

	resp, err := decodeResponse(conn)
	if err != nil {
		return "", fmt.Errorf("falha ao receber resposta AUTH: %w", err)
	}

	if ok := resp.GetOk(); ok != nil {
		if ok.Nome != "" {
			token := strings.TrimSpace(ok.Nome)
			fmt.Printf(" <- Autenticação OK. Nome: %s\n", token)
			return token, nil
		}
	} else if erro := resp.GetErro(); erro != nil {
		return "", fmt.Errorf("autenticação falhou: %s", erro.Mensagem)
	}

	return "", fmt.Errorf("resposta de autenticação inválida ou token vazio")
}

func RunOperation(conn net.Conn, opName string, token string, args map[string]string) {
	fmt.Printf("\n--- 2. OPERAÇÃO: %s ---\n", strings.ToUpper(opName))

	protobufParams := make([]string, 0)

	switch opName {
	case "echo":
		if val, ok := args["mensagem"]; ok {
			protobufParams = append(protobufParams, fmt.Sprintf("mensagem=%s", val))
		}
	case "soma":
		if val, ok := args["numeros"]; ok {
			protobufParams = append(protobufParams, fmt.Sprintf("numeros=%s", val))
		}
	case "status":
		if val, ok := args["detalhado"]; ok {
			protobufParams = append(protobufParams, fmt.Sprintf("detalhado=%s", val))
		}
	case "historico":
		if val, ok := args["limite"]; ok {
			protobufParams = append(protobufParams, fmt.Sprintf("limite=%s", val))
		}
	case "timestamp":
	}

	operacaoReq := &pb.Requisicao{
		Tipo: &pb.Requisicao_Operacao{
			Operacao: &pb.Operacao{
				Token:        token,
				NomeOperacao: opName,
				Parametros:   protobufParams,
			},
		},
	}

	encodedMsg, err := encodeMessage(operacaoReq)
	if err != nil {
		fmt.Printf("⚠️ Operação %s Falhou (Encodar): %v\n", opName, err)
		return
	}

	if _, err := conn.Write(encodedMsg); err != nil {
		fmt.Printf("⚠️ Operação %s Falhou (Enviar): %v\n", opName, err)
		return
	}

	paramsLog := strings.Join(protobufParams, ", ")
	fmt.Printf(" -> Enviando %d bytes (Op: %s, Params: [%s])...\n", len(encodedMsg), opName, paramsLog)

	resp, err := decodeResponse(conn)
	if err != nil {
		fmt.Printf("⚠️ Operação %s Falhou (Receber): %v\n", opName, err)
		return
	}

	if ok := resp.GetOk(); ok != nil {
		fmt.Printf(" ✅ Operação %s Sucesso! Resposta no token/nome/matricula.\n", opName)
		if ok.Nome != "" {
			fmt.Printf("    -> Nome/Retorno: %s\n", ok.Nome)
		}
	} else if erro := resp.GetErro(); erro != nil {
		fmt.Printf(" ⚠️ Operação %s Falhou: %s\n", opName, erro.Mensagem)
	} else {
		fmt.Printf(" ⚠️ Operação %s Falhou: Resposta desconhecida do servidor.\n", opName)
	}
}

func Logout(conn net.Conn, token string) {
	fmt.Println("\n--- 3. LOGOUT ---")

	logoutReq := &pb.Requisicao{
		Tipo: &pb.Requisicao_Logout{
			Logout: &pb.Logout{
				Token: token,
			},
		},
	}

	encodedMsg, _ := encodeMessage(logoutReq)

	if _, err := conn.Write(encodedMsg); err != nil {
		fmt.Printf("⚠️ Aviso: Falha ao enviar Logout: %v\n", err)
	}

	resp, err := decodeResponse(conn)
	if err != nil {
		fmt.Printf("⚠️ Aviso: Falha ao receber confirmação de Logout: %v\n", err)
	} else if ok := resp.GetOk(); ok != nil {
		fmt.Println("✅ Logout confirmado pelo servidor.")
	} else if erro := resp.GetErro(); erro != nil {
		fmt.Printf("⚠️ Logout falhou: %s\n", erro.Mensagem)
	}

	conn.Close()
	fmt.Println("✅ Conexão TCP local fechada.")
}

func main() {
	fmt.Printf("--- Cliente Protobuf Protocol (%s) ---\n", SERVER_ADDR)

	conn, err := net.Dial(PROTOCOL, SERVER_ADDR)
	if err != nil {
		fmt.Printf("❌ Erro ao conectar a %s: %v\n", SERVER_ADDR, err)
		os.Exit(1)
	}

	fmt.Println("✅ Conexão TCP estabelecida.")

	sessionToken, err := Authenticate(conn, ALUNO_ID)
	if err != nil {
		fmt.Printf("❌ Erro de Autenticação: %v\n", err)
		conn.Close()
		return
	}

	tokenPreview := sessionToken
	if len(sessionToken) > 10 {
		tokenPreview = sessionToken[:10]
	}
	fmt.Printf("✅ Autenticação bem-sucedida! Token: %s...\n", tokenPreview)

	RunOperation(conn, "echo", sessionToken, map[string]string{
		"mensagem": "Mensagem de teste binária",
	})

	RunOperation(conn, "soma", sessionToken, map[string]string{
		"numeros": "10,20,30",
	})

	RunOperation(conn, "status", sessionToken, map[string]string{
		"detalhado": "true",
	})

	RunOperation(conn, "historico", sessionToken, map[string]string{
		"limite": "2",
	})

	RunOperation(conn, "timestamp", sessionToken, map[string]string{})

	Logout(conn, sessionToken)

	fmt.Println("\n--- Cliente Protobuf: Fim das operações. ---")
}
