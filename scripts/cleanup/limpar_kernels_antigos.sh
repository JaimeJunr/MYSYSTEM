#!/usr/bin/env bash
# Remove kernels antigos em sistemas apt (autoremove). Requer reinícios após atualizações.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Kernels antigos (apt) ---"
echo "Kernel em execução: $(uname -r)"
echo ""

if ! command -v apt &>/dev/null; then
  echo "apt não encontrado — use a ferramenta da sua distro (dnf/zypper/pacman) manualmente."
  exit 0
fi

confirm_action "Apt autoremove/autoclean (remove kernels e headers antigos já marcados como órfãos)" \
  "apt-get autoremove --purge -y && apt-get autoclean -y" \
  "dpkg-query -W -f='\${Package}\n' 2>/dev/null | grep -c '^linux-image-' || echo 0"

echo "--- Concluído ---"
