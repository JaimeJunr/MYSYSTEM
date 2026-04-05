#!/bin/bash

# Script de Limpeza Geral (Caches e Sistema)
# Parte do sistema de limpeza SOLID

REAL_USER=${REAL_USER:-${SUDO_USER:-$USER}}
REAL_HOME=${REAL_HOME:-$(getent passwd "$REAL_USER" | cut -d: -f6)}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Iniciando Limpeza Geral ---"

# 1. Lixeira
if [ -d "$REAL_HOME/.local/share/Trash" ]; then
    confirm_action "1. Limpar Lixeira" \
        "rm -rf \"$REAL_HOME/.local/share/Trash\"/*" \
        "du -sh \"$REAL_HOME/.local/share/Trash\" 2>/dev/null | cut -f1"
else
    echo "1. Lixeira (vazia/não encontrada). Pular."
    echo ""
fi

# 2. Docker
if command -v docker &> /dev/null; then
    confirm_action "2. Limpar Docker (imagens não usadas + build cache)" \
        "docker system prune -af --volumes" \
        "docker system df 2>/dev/null | grep 'Total' -A 1 | tail -n 1 | awk '{print \$4}'"
else
    echo "2. Docker não instalado. Pular."
    echo ""
fi

# 3. Poetry
if [ -d "$REAL_HOME/.cache/pypoetry" ]; then
    confirm_action "3. Limpar cache do Poetry" \
        "sudo -u \"$REAL_USER\" poetry cache clear pypi --all 2>/dev/null || rm -rf \"$REAL_HOME/.cache/pypoetry\"/*" \
        "du -sh \"$REAL_HOME/.cache/pypoetry\" 2>/dev/null | cut -f1"
fi

# 4. Pip
if [ -d "$REAL_HOME/.cache/pip" ]; then
    confirm_action "4. Limpar cache do pip" \
        "sudo -u \"$REAL_USER\" pip cache purge 2>/dev/null || rm -rf \"$REAL_HOME/.cache/pip\"/*" \
        "du -sh \"$REAL_HOME/.cache/pip\" 2>/dev/null | cut -f1"
fi

# 5. Systemd Logs
if command -v journalctl &> /dev/null; then
    confirm_action "5. Limpar logs do systemd (mantém últimos 3 dias - Mais Intenso)" \
        "sudo journalctl --vacuum-time=3d" \
        "journalctl --disk-usage 2>/dev/null | grep -oP '\d+\.\d+G'"
fi

# 6. Chrome (Intense)
if [ -d "$REAL_HOME/.cache/google-chrome" ]; then
    confirm_action "6. Limpar cache do Chrome (AVISO: Pode deslogar/lento)" \
        "rm -rf \"$REAL_HOME/.cache/google-chrome\"/*" \
        "du -sh \"$REAL_HOME/.cache/google-chrome\" 2>/dev/null | cut -f1"
fi

# 7. Thumbnails (New)
if [ -d "$REAL_HOME/.cache/thumbnails" ]; then
    confirm_action "7. Limpar Thumbnails (Miniaturas de arquivos)" \
        "rm -rf \"$REAL_HOME/.cache/thumbnails\"/*" \
        "du -sh \"$REAL_HOME/.cache/thumbnails\" 2>/dev/null | cut -f1"
fi

# 8. Apt Autoremove (Remove dependências não utilizadas + clean)
if command -v apt &> /dev/null; then
    confirm_action "8. Apt Autoremove (Remove dependências não utilizadas + clean)" \
        "sudo apt autoremove -y && sudo apt clean" \
        "sudo du -sh /var/cache/apt/archives 2>/dev/null | cut -f1"
fi

# 9. Snap Disabled
if command -v snap &> /dev/null; then
    echo "---------------------------------------------------"
    echo "9. Verificando Snaps desabilitados..."
    DISABLED_SNAPS=$(snap list --all 2>/dev/null | grep disabled | awk '{print $1, $3}' | sort -u -k1,1)

    if [[ "${HOMESTEAD_DRY_RUN:-}" == "1" ]]; then
        if [ -n "$DISABLED_SNAPS" ]; then
            echo "   [DRY-RUN] Revisões desativadas que seriam removíveis:"
            echo "$DISABLED_SNAPS"
            echo "   [DRY-RUN] Comando por linha: sudo snap remove NAME --revision=REV"
        else
            echo "   [DRY-RUN] Nenhum snap desabilitado listado."
        fi
        echo ""
    elif [ -n "$DISABLED_SNAPS" ]; then
        echo "   Snaps desabilitados encontrados:"
        echo "$DISABLED_SNAPS"
        echo ""
        read -p "   Deseja remover estes snaps desabilitados? [s/N]: " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Ss]$ ]]; then
            echo "$DISABLED_SNAPS" | while read -r snap_name revision; do
                echo "   Removendo $snap_name (revisão $revision)..."
                sudo snap remove "$snap_name" --revision="$revision" 2>/dev/null || true
            done
            echo "   ✓ Snaps desabilitados removidos."
        else
            echo "   ⊘ Pular."
        fi
    else
        echo "   Nenhum snap desabilitado encontrado."
    fi
    echo ""
fi

# 10. Yarn
if [ -d "$REAL_HOME/.cache/yarn" ]; then
    confirm_action "10. Limpar cache do Yarn" \
        "sudo -u \"$REAL_USER\" yarn cache clean 2>/dev/null || rm -rf \"$REAL_HOME/.cache/yarn\"/*" \
        "du -sh \"$REAL_HOME/.cache/yarn\" 2>/dev/null | cut -f1"
fi

# 11. NPM
if [ -d "$REAL_HOME/.npm" ]; then
    confirm_action "11. Limpar cache do NPM" \
        "sudo -u \"$REAL_USER\" npm cache clean --force 2>/dev/null || rm -rf \"$REAL_HOME/.npm\"/*" \
        "du -sh \"$REAL_HOME/.npm\" 2>/dev/null | cut -f1"
fi

# 12. Gradle
if [ -d "$REAL_HOME/.gradle/caches" ]; then
    confirm_action "12. Limpar cache do Gradle (AVISO: Re-download pesado)" \
        "rm -rf \"$REAL_HOME/.gradle/caches\"/*" \
        "du -sh \"$REAL_HOME/.gradle/caches\" 2>/dev/null | cut -f1"
fi

# 13. Cargo
if [ -d "$REAL_HOME/.cargo/registry" ]; then
    confirm_action "13. Limpar cache do Cargo/Rust (registry)" \
        "rm -rf \"$REAL_HOME/.cargo/registry\"/*" \
        "du -sh \"$REAL_HOME/.cargo/registry\" 2>/dev/null | cut -f1"
fi

# 14. Maven
if [ -d "$REAL_HOME/.m2/repository" ]; then
    confirm_action "14. Limpar repositório Maven (AVISO: Re-download pesado)" \
        "rm -rf \"$REAL_HOME/.m2/repository\"/*" \
        "du -sh \"$REAL_HOME/.m2/repository\" 2>/dev/null | cut -f1"
fi
