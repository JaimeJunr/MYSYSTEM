#!/usr/bin/env bash
# Flatpak (unused) e revisões snap desativadas.

REAL_USER=${REAL_USER:-${SUDO_USER:-$USER}}
REAL_HOME=${REAL_HOME:-$(getent passwd "$REAL_USER" | cut -d: -f6)}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Limpeza Flatpak / Snap ---"

if command -v flatpak &>/dev/null; then
  confirm_action "Flatpak: remover runtimes e apps não usados (flatpak uninstall --unused)" \
    "flatpak uninstall --unused -y 2>/dev/null || flatpak uninstall --unused" \
    "du -sh \"$REAL_HOME/.local/share/flatpak\" 2>/dev/null | cut -f1"
else
  echo "Flatpak não instalado — a saltar."
fi

if command -v snap &>/dev/null; then
  echo "---------------------------------------------------"
  echo "Snaps desabilitados..."
  DISABLED_SNAPS=$(snap list --all 2>/dev/null | grep disabled | awk '{print $1, $3}' | sort -u -k1,1)
  if [[ "${HOMESTEAD_DRY_RUN:-}" == "1" ]]; then
    if [ -n "$DISABLED_SNAPS" ]; then
      echo "[DRY-RUN] Revisões a remover com confirmação normal:"
      echo "$DISABLED_SNAPS"
    else
      echo "[DRY-RUN] Nenhuma revisão disabled."
    fi
  elif [ -n "$DISABLED_SNAPS" ]; then
    echo "$DISABLED_SNAPS"
    read -r -p "Remover estas revisões? [s/N]: " -n 1
    echo ""
    if [[ $REPLY =~ ^[Ss]$ ]]; then
      echo "$DISABLED_SNAPS" | while read -r snap_name revision; do
        echo "Removendo $snap_name ($revision)..."
        sudo snap remove "$snap_name" --revision="$revision" 2>/dev/null || true
      done
    fi
  else
    echo "Nenhuma revisão disabled."
  fi
else
  echo "Snap não instalado — a saltar."
fi

echo "--- Concluído ---"
