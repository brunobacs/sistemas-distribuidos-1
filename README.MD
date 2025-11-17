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
  * **Implementação:** O cliente (`string.go`) usa a função `formatStringsMessage` para construir o payload exato, incluindo os pares **chave=valor** e o terminador **`|FIM\n`**, que não estava na documentação original.
  * **Desafio:** **Reversão de Engenharia do Terminador.** O servidor exigia o terminador não-ppadrão `|FIM` e o formato `chave=valor` em todas as requisições (`AUTH`, `OP`), o que exigiu múltiplas iterações para depurar.

### 2\. Cliente Protocolo JSON (Porta 8081)

  * **Protocolo:** Comunicação baseada em objetos JSON (`\n`)
  * **Formato Descoberto:** A estrutura real não era a da documentação original, mas sim **`{"tipo": "...", "operacao": "...", "parametros": {...}}`**.
  * **Implementação:** O cliente (`Json/main.go`) utiliza structs dedicadas (`AuthRequest`, `OperationRequest`) com tags JSON minúsculas (`json:"aluno_id"`, `json:"parametros"`) para garantir a serialização correta.
  * **Desafio:** **Falha de Servidor e Timeout.** O servidor JSON falhou consistentemente com `i/o timeout` na leitura da resposta (mesmo com requisições corretas), indicando um provável **erro de tratamento de JSON** ou **falha de firewall de saída** no lado do servidor. O cliente foi implementado com a estrutura correta, mas a comunicação falha devido a problemas externos.

### 3\. Cliente Protocolo Protocol Buffers (Porta 8082)

  * **Protocolo:** Comunicação binária com cabeçalho de 4 bytes (tamanho)
  * **Formato Descoberto:** A autenticação e as operações exigem a inclusão do campo **`timestamp`** (formato RFC3339Nano) e o uso de um **`map<string, string>`** para parâmetros na mensagem `Operacao`.
  * **Implementação:** O cliente (`ProtocolBuffer/protocolbuff-client.go`) implementa:
      * **Codificação/Decodificação Binária:** Funções `encodeMessage` (cria o cabeçalho de 4 bytes + payload) e `decodeResponse` (lê o cabeçalho e o payload).
      * **Estrutura Correta:** Uso do `oneof conteudo` e `map[string]string Parametros` (o formato real do servidor) para evitar o erro de serialização binária.
  * **Desafio:** **Incompatibilidade de Schema Binário.** Foi necessário inferir e aplicar o schema *real* do servidor (incluindo `Timestamp` e o `map<string, string>`) em vez do `repeated string` inicial, para que o servidor aceitasse o payload.

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

**Nota sobre JSON (Opção 2):** Se o cliente JSON ainda retornar `i/o timeout`, isso confirma a falha no lado do servidor, conforme documentado na seção Desafios Críticos.