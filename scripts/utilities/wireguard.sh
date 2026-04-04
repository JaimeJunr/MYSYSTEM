#!/usr/bin/env bash
# WireGuard (ferramentas de linha de comando / módulo conforme distro).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando componentes WireGuard…"
if hs_is_debian_like; then
  hs_install_packages wireguard
elif hs_is_arch_like || [[ "$ID" == "solus" ]]; then
  hs_install_packages wireguard-tools
elif hs_is_fedora_like || hs_is_suse; then
  hs_install_packages wireguard-tools
else
  hs_die "Distribuição não suportada para instalação automática (${ID})."
fi
hs_info "WireGuard instalado."
