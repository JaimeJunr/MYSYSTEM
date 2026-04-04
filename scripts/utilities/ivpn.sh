#!/usr/bin/env bash
# IVPN (repositório oficial).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando IVPN…"
if [[ "$ID" == "debian" ]] || [[ "${ID_LIKE:-}" == *debian* ]]; then
  curl -fsSL https://repo.ivpn.net/stable/debian/generic.gpg | sudo gpg --dearmor -o /usr/share/keyrings/ivpn-archive-keyring.gpg
  curl -fsSL https://repo.ivpn.net/stable/debian/generic.list | sudo tee /etc/apt/sources.list.d/ivpn.list
  sudo apt-get update
  hs_install_packages ivpn-ui
elif [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
  curl -fsSL https://repo.ivpn.net/stable/ubuntu/generic.gpg | sudo gpg --dearmor -o /usr/share/keyrings/ivpn-archive-keyring.gpg
  curl -fsSL https://repo.ivpn.net/stable/ubuntu/generic.list | sudo tee /etc/apt/sources.list.d/ivpn.list
  sudo apt-get update
  hs_install_packages ivpn-ui
elif hs_is_fedora_like; then
  if hs_has_ostree; then
    curl -fsSLo /tmp/ivpn.repo https://repo.ivpn.net/stable/fedora/generic/ivpn.repo
    sudo install -o root -g root -m644 /tmp/ivpn.repo /etc/yum.repos.d/ivpn.repo
    rm -f /tmp/ivpn.repo
    sudo rpm-ostree refresh-md || true
  else
    sudo dnf config-manager addrepo --from-repofile=https://repo.ivpn.net/stable/fedora/generic/ivpn.repo
  fi
  hs_install_packages ivpn-ui
else
  hs_die "Distribuição não suportada (${ID})."
fi
hs_info "Concluído."
