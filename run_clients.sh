#!/bin/bash


CLIENT_STRING_DIR="Socket/string"
CLIENT_JSON_DIR="Json"
CLIENT_PB_DIR="ProtocolBuffer"

run_go_client() {
    local dir="$1"
    
    (
        cd "$dir" || { echo "Erro: Diretório $dir não encontrado."; exit 1; }
        
        echo -e "\n========================================================"
        echo "EXECUTANDO CLIENTE EM: $(pwd)"
        echo "========================================================"
    
        
        if [ "$dir" = "$CLIENT_STRING_DIR" ]; then
            go run string.go
        elif [ "$dir" = "$CLIENT_JSON_DIR" ]; then
            go run .
        elif [ "$dir" = "$CLIENT_PB_DIR" ]; then
            go run .
        fi
        
        if [ $? -eq 0 ]; then
            echo -e "\n✅ EXECUÇÃO CONCLUÍDA COM SUCESSO."
        else
            echo -e "\n❌ EXECUÇÃO FALHOU. Verifique os erros de compilação/runtime."
        fi
    )
}


while true; do
    echo -e "\n========================================================"
    echo "MENU DE CLIENTES DE SISTEMAS DISTRIBUÍDOS"
    echo "========================================================"
    echo "1) ⚡ Executar Cliente String (Porta 8080)"
    echo "2) 📝 Executar Cliente JSON (Porta 8081)"
    echo "3) 📦 Executar Cliente Protocol Buffer (Porta 8082)"
    echo "4) 🚪 Sair"
    echo "--------------------------------------------------------"
    read -p "Escolha uma opção: " choice
    
    case $choice in
        1)
            run_go_client "$CLIENT_STRING_DIR"
            ;;
        2)
            run_go_client "$CLIENT_JSON_DIR"
            ;;
        3)
            run_go_client "$CLIENT_PB_DIR"
            ;;
        4)
            echo -e "\nSaindo. Obrigado! 👋"
            exit 0
            ;;
        *)
            echo -e "\nOpção inválida. Por favor, escolha 1, 2, 3 ou 4."
            ;;
    esac
done