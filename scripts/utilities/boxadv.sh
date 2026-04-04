#!/usr/bin/env bash
# Distrobox-Adv BR + pcscd + DistroShelf (certificado digital — advogados).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando dependências (distrobox, podman, smartcard)…"
if hs_is_fedora_like; then
  hs_install_packages distrobox podman pcsc-lite pcsc-lite-ccid
elif [[ "$ID" == "debian" || "$ID" == "ubuntu" ]] || hs_is_debian_like; then
  if [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
    sudo add-apt-repository -y ppa:michel-slm/distrobox
    sudo apt-get update
  fi
  hs_install_packages distrobox podman pcsc-lite ccid
elif hs_is_arch_like; then
  hs_install_packages distrobox podman pcsclite ccid
elif hs_is_suse; then
  hs_install_packages distrobox podman pcsc-ccid
else
  hs_die "Distribuição não mapeada (${ID})."
fi
sudo systemctl enable --now pcscd.service
if [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
  distrobox-assemble create --file https://raw.githubusercontent.com/pedrohqb/distrobox-adv-br/refs/heads/main/distrobox-adv-br-legado
else
  distrobox-assemble create --file https://raw.githubusercontent.com/pedrohqb/distrobox-adv-br/refs/heads/main/distrobox-adv-br
fi
hs_flatpak_install_user com.ranfdev.DistroShelf
hs_info "Concluído."
