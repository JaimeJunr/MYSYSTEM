#!/usr/bin/env bash
# Reduzir tamanho do journal systemd; lista maiores em /var/log (só leitura).

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Journal e logs do sistema ---"

if command -v journalctl &>/dev/null; then
  confirm_action "Journal: manter apenas últimos 7 dias (vacuum-time)" \
    "journalctl --vacuum-time=7d" \
    "journalctl --disk-usage 2>/dev/null | head -1"
else
  echo "journalctl não disponível."
fi

if [[ -d /var/log ]]; then
  echo ""
  echo "Maiores entradas em /var/log (só leitura):"
  du -h /var/log 2>/dev/null | sort -hr | head -20 || true
fi

echo ""
echo "--- Concluído ---"
