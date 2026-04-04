#!/usr/bin/env bash
# Input Remapper (input-remapper).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando Input Remapper…"
if [[ "$ID" == "arch" || "$ID" == "cachyos" ]] || hs_is_arch_like; then
  hs_chaotic_aur_enable
  hs_install_packages input-remapper-git
else
  hs_install_packages input-remapper
fi
hs_info "Pode ser necessário reiniciar o sistema."
