#!/usr/bin/env bash
# NordVPN (script oficial ou AUR em Arch).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando NordVPN…"
if hs_is_debian_like || hs_is_fedora_like || hs_is_suse; then
  sh <(wget -qO- https://downloads.nordcdn.com/apps/linux/install.sh) -p nordvpn-gui -n
elif hs_is_arch_like; then
  hs_chaotic_aur_enable
  hs_install_packages nordvpn-bin
else
  hs_die "Distribuição não suportada (${ID})."
fi
hs_info "Concluído."
