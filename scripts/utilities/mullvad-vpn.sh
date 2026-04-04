#!/usr/bin/env bash
# Mullvad VPN (repositório oficial).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Configurando repositório Mullvad…"
if hs_is_debian_like; then
  curl -fsSLo /tmp/mullvad-keyring.asc https://repository.mullvad.net/deb/mullvad-keyring.asc
  sudo install -o root -g root -m644 /tmp/mullvad-keyring.asc /usr/share/keyrings/mullvad-keyring.asc
  rm -f /tmp/mullvad-keyring.asc
  echo "deb [signed-by=/usr/share/keyrings/mullvad-keyring.asc arch=$(dpkg --print-architecture)] https://repository.mullvad.net/deb/stable stable main" | sudo tee /etc/apt/sources.list.d/mullvad.list
  sudo apt-get update
  hs_install_packages mullvad-vpn
elif hs_is_fedora_like; then
  if hs_has_ostree; then
    curl -fsSLo /tmp/mullvad.repo https://repository.mullvad.net/rpm/stable/mullvad.repo
    sudo install -o root -g root -m644 /tmp/mullvad.repo /etc/yum.repos.d/mullvad.repo
    rm -f /tmp/mullvad.repo
    sudo rpm-ostree refresh-md || true
  else
    sudo dnf config-manager addrepo --from-repofile=https://repository.mullvad.net/rpm/stable/mullvad.repo
  fi
  hs_install_packages mullvad-vpn
elif hs_is_arch_like; then
  hs_chaotic_aur_enable
  hs_install_packages mullvad-vpn
else
  hs_die "Distribuição não suportada para instalação automática (${ID})."
fi
hs_info "Mullvad instalado."
