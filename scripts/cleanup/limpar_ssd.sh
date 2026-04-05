#!/bin/bash

# Script Orquestrador de Limpeza de SSD
# Aplica o Princípio da Responsabilidade Única (SOLID)

set -e

REAL_USER=${SUDO_USER:-$USER}
REAL_HOME=$(getent passwd "$REAL_USER" | cut -d: -f6)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "=== Limpeza de Espaço no SSD (Orquestrador) ==="
echo "Usuário: $REAL_USER"
echo "Home: $REAL_HOME"
echo ""

show_space() {
    echo "Espaço disponível:"
    df -h / | tail -1
    echo ""
}

export REAL_USER REAL_HOME
export -f confirm_action

# --- Execução ---
show_space

# 1. Limpeza Geral (Caches, Sistema, etc)
if [ -f "$SCRIPT_DIR/limpar_geral.sh" ]; then
    source "$SCRIPT_DIR/limpar_geral.sh"
else
    echo "⚠ Erro: limpar_geral.sh não encontrado."
fi

# 2. Limpeza de Itens Grandes (Arquivos e Pastas)
if [ -f "$SCRIPT_DIR/limpar_grandes.sh" ]; then
    source "$SCRIPT_DIR/limpar_grandes.sh"
else
    echo "⚠ Erro: limpar_grandes.sh não encontrado."
fi

show_space
echo "=== Limpeza Concluída! ==="
