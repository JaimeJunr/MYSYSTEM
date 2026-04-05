#!/usr/bin/env bash
# Miniaturas em ~/.cache/thumbnails

REAL_USER=${REAL_USER:-${SUDO_USER:-$USER}}
REAL_HOME=${REAL_HOME:-$(getent passwd "$REAL_USER" | cut -d: -f6)}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Thumbnails (miniaturas) ---"

if [ -d "$REAL_HOME/.cache/thumbnails" ]; then
  confirm_action "Apagar miniaturas em ~/.cache/thumbnails (regeneram ao navegar)" \
    "rm -rf \"$REAL_HOME/.cache/thumbnails\"/*" \
    "du -sh \"$REAL_HOME/.cache/thumbnails\" 2>/dev/null | cut -f1"
else
  echo "Pasta ~/.cache/thumbnails não existe."
fi

echo "--- Concluído ---"
