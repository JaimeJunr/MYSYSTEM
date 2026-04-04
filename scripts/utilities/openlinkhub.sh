#!/usr/bin/env bash
# OpenLinkHub — periféricos Corsair (repositórios oficiais / COPR / AUR bin).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
tag="$(curl -fsSL "https://api.github.com/repos/jurkovic-nikola/OpenLinkHub/releases/latest" | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p')"
[[ -n "$tag" ]] || hs_die "Não foi possível obter a versão mais recente do OpenLinkHub."
hs_info "Instalando OpenLinkHub (${tag})…"
if [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
  sudo add-apt-repository -y ppa:jurkovic-nikola/openlinkhub
  sudo apt-get update
  sudo apt-get install -y openlinkhub
elif [[ "$ID" == "debian" ]]; then
  cd "${TMPDIR:-/tmp}"
  wget "https://github.com/jurkovic-nikola/OpenLinkHub/releases/download/${tag}/OpenLinkHub_${tag}_amd64.deb"
  sudo apt-get install -y "./OpenLinkHub_${tag}_amd64.deb"
  rm -f "./OpenLinkHub_${tag}_amd64.deb"
elif hs_has_ostree; then
  cd "${TMPDIR:-/tmp}"
  fv="$(rpm -E %fedora)"
  wget "https://copr.fedorainfracloud.org/coprs/jurkovic-nikola/OpenLinkHub/repo/fedora-${fv}/jurkovic-nikola-OpenLinkHub-fedora-${fv}.repo"
  sudo install -o 0 -g 0 "jurkovic-nikola-OpenLinkHub-fedora-${fv}.repo" "/etc/yum.repos.d/"
  rm -f "jurkovic-nikola-OpenLinkHub-fedora-${fv}.repo"
  sudo rpm-ostree refresh-md || true
  sudo rpm-ostree install -yA OpenLinkHub
elif hs_is_fedora_like; then
  sudo dnf copr enable -y jurkovic-nikola/OpenLinkHub
  sudo dnf install -y OpenLinkHub
elif hs_is_arch_like; then
  hs_install_packages openlinkhub-bin
else
  hs_die "Distribuição não suportada (${ID})."
fi
sudo systemctl enable --now OpenLinkHub.service || true
sleep 1
xdg-open http://127.0.0.1:27003 2>/dev/null || true
hs_info "OpenLinkHub configurado."
