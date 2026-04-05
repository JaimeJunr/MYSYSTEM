#!/usr/bin/env bash
# Funções partilhadas por scripts de limpeza Homestead.
# HOMESTEAD_DRY_RUN=1: apenas descreve ações e tamanhos, sem alterar o sistema.

homestead_is_dry_run() {
  [[ "${HOMESTEAD_DRY_RUN:-}" == "1" ]]
}

# Uso: confirm_action "título" "comando" "cmd opcional para du/tamanho"
confirm_action() {
  local description="$1"
  local command_to_run="$2"
  local size_check_cmd="${3:-}"

  echo "---------------------------------------------------"
  echo "$description"

  if [[ -n "$size_check_cmd" ]]; then
    local size
    size=$(eval "$size_check_cmd" 2>/dev/null || echo "N/A")
    echo "   Tamanho estimado: $size"
  fi

  if homestead_is_dry_run; then
    echo "   [DRY-RUN] Comando que seria executado:"
    echo "   $command_to_run"
    echo ""
    return 0
  fi

  read -r -p "   Deseja prosseguir? [s/N]: " -n 1
  echo ""
  if [[ $REPLY =~ ^[Ss]$ ]]; then
    eval "$command_to_run"
    echo "   ✓ Concluído."
  else
    echo "   ⊘ Pular."
  fi
  echo ""
}
