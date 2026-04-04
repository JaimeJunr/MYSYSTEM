#!/usr/bin/env bash
# Podman + Distrobox + DistroShelf (Flatpak).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando podman, distrobox e DistroShelf…"
if hs_is_fedora_like; then
  hs_install_packages podman distrobox
elif [[ "$ID" == "debian" || "$ID" == "ubuntu" ]] || hs_is_debian_like; then
  hs_install_packages podman distrobox
elif hs_is_arch_like; then
  hs_install_packages podman distrobox
elif hs_is_suse; then
  hs_install_packages podman distrobox
else
  hs_die "Distribuição não mapeada para podman/distrobox (${ID})."
fi
hs_flatpak_install_user com.ranfdev.DistroShelf
hs_info "Concluído."
