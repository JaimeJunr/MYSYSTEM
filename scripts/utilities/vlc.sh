#!/usr/bin/env bash
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando VLC…"
if hs_is_fedora_like && ! hs_has_ostree; then
  hs_rpmfusion_enable
  hs_install_packages vlc libavcodec-freeworld
elif hs_is_suse; then
  hs_install_packages opi
  sudo opi -n codecs || true
  hs_install_packages vlc
else
  hs_install_packages vlc
fi
hs_info "Concluído."
