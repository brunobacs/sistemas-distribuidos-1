# 📚 Clientes para Protocolos de Comunicação

## Visão Geral

Este projeto foi desenvolvido para a disciplina de **Sistemas Distribuídos** com o objetivo de implementar e comparar clientes funcionais em Go para três protocolos de comunicação distintos que interagem com servidores remotos na AWS.

O principal foco foi garantir a correta serialização/desserialização e superar desafios de implementação de protocolos não-padrão.

-----

## 🏗️ Estrutura do Projeto

O projeto é organizado modularmente, onde cada cliente reside em seu próprio diretório e tem seu próprio arquivo `go.mod` para gerenciamento de dependências.

```
/projeto-sd/
├── Json/                     # Cliente Protocolo JSON (Porta 8081)
│   ├── main.go               # Lógica principal, funções de operação
│   └── types.go              # Definições das Structs JSON (DTOs)
├── ProtocolBuffer/           # Cliente Protocolo Protobuf (Porta 8082)
│   ├── client-api.proto      # Schema Protobuf
│   ├── client-api.pb.go      # Código Go Gerado
│   └── protocolbuff-client.go # Lógica principal e I/O binário
└── Socket/
    └── string/               # Cliente Protocolo Strings (Porta 8080)
        └── string.go           # Lógica principal (Comunicação por pipe/linha)
```

-----

## ⚙️ Implementações e Desafios Críticos

### 1\. Cliente Protocolo Strings (Porta 8080)

  * **Protocolo:** Comunicação baseada em linhas de texto (`|`)
  * **Formato Descoberto:** `COMANDO|chave1=valor1|chave2=valor2|...|FIM\n`
  * **Implementação:** O cliente (`string.go`) usa a função `formatStringsMessage` para construir o payload exato, incluindo os pares **chave=valor** e o terminador **`|FIM\n`**.
  * **Desafio:** **Reversão de Engenharia do Terminador.** O servidor exigia o terminador `|FIM` e o formato `chave=valor` em todas as requisições (`AUTH`, `OP`), o que exigiu múltiplas iterações para depurar.

### 2\. Cliente Protocolo JSON (Porta 8081)

  * **Protocolo:** Comunicação baseada em objetos JSON (`\n`)
  * **Formato Descoberto:** A estrutura real era **`{"tipo": "...", "operacao": "...", "parametros": {...}}`**.
  * **Implementação:** O cliente (`Json/main.go`) utiliza structs dedicadas (`AuthRequest`, `OperationRequest`) com tags JSON minúsculas (`json:"aluno_id"`, `json:"parametros"`) para garantir a serialização correta.


### 3\. Cliente Protocolo Protocol Buffers (Porta 8082)

  * **Protocolo:** Comunicação binária com cabeçalho de 4 bytes (tamanho)
  * **Formato Descoberto:** A autenticação e as operações exigem a inclusão do campo **`timestamp`** e o uso de um **`map<string, string>`** para parâmetros na mensagem `Operacao`.

-----

## 🚀 Como Executar os Clientes

O projeto inclui um script shell na raiz para facilitar a execução e o teste de cada cliente.

### Requisitos

  * Go 1.21+ instalado e configurado (`$GOPATH/bin` no `$PATH`).
  * O compilador `protoc` instalado para o cliente Protobuf.

### 1\. Dar Permissão e Executar o Script

Na raiz do projeto, execute:

```bash
chmod +x run_clients.sh
./run_clients.sh
```

### 2\. Menu de Opções

O script apresentará as seguintes opções para testar cada protocolo individualmente:

```
========================================================
MENU DE CLIENTES DE SISTEMAS DISTRIBUÍDOS
========================================================
1) ⚡ Executar Cliente String (Porta 8080)
2) 📝 Executar Cliente JSON (Porta 8081)
3) 📦 Executar Cliente Protocol Buffer (Porta 8082)
4) 🚪 Sair
```
🎥 Apresentação
Para mais informações sobre a comparação de desempenho e eficiência de Protocol Buffers versus JSON, consulte o vídeo:

YouTube: Protocol Buffers vs JSON vs XML vs CSV Performance Comparison