#!/usr/bin/env bash
# Cloudflare WARP client (repositório oficial).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
if command -v warp-cli &>/dev/null; then
  hs_info "warp-cli já está instalado."
  exit 0
fi
hs_info "Instalando Cloudflare WARP…"
if command -v apt-get &>/dev/null; then
  curl -fsSL https://pkg.cloudflareclient.com/pubkey.gpg | sudo gpg --yes --dearmor --output /usr/share/keyrings/cloudflare-warp-archive-keyring.gpg
  rel="$(grep VERSION_CODENAME /etc/os-release | cut -d= -f2)"
  supported="noble jammy focal bookworm bullseye trixie"
  if [[ ! " $supported " == *" $rel "* ]]; then
    rel="jammy"
    hs_warn "Codename ${rel} não listado como suportado; usando fallback jammy no repositório Cloudflare."
  fi
  echo "deb [signed-by=/usr/share/keyrings/cloudflare-warp-archive-keyring.gpg] https://pkg.cloudflareclient.com/ ${rel} main" | sudo tee /etc/apt/sources.list.d/cloudflare-client.list
  sudo apt-get update
  sudo apt-get install -y cloudflare-warp
elif command -v dnf &>/dev/null || command -v yum &>/dev/null; then
  curl -fsSL https://pkg.cloudflareclient.com/cloudflare-warp-ascii.repo | sudo tee /etc/yum.repos.d/cloudflare-warp.repo
  hs_install_packages cloudflare-warp
else
  hs_die "Use Debian/Ubuntu ou Fedora/RHEL para instalação automática do WARP."
fi
hs_info "Concluído. Configure com: warp-cli register"
