#!/usr/bin/env bash
# Piper (ratbagd + Flatpak sistema).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando dependências Piper…"
if [[ "$ID" == "ubuntu" || "$ID" == "debian" || "${ID_LIKE:-}" == *ubuntu* ]]; then
  hs_install_packages ratbagd
elif hs_is_arch_like || [[ "$ID" == "solus" ]]; then
  hs_install_packages libratbag
elif hs_is_fedora_like || hs_has_ostree; then
  hs_install_packages libratbag-ratbagd
elif hs_is_suse; then
  hs_install_packages libratbag-tools
else
  hs_warn "Pacote ratbag não mapeado; tentando ratbagd…"
  hs_install_packages ratbagd || true
fi
hs_flatpak_install_system org.freedesktop.Piper
hs_info "Piper instalado. Pode ser necessário reiniciar."
