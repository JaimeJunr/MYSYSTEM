#!/usr/bin/env bash
# Logitech Unifying / Bolt — Solaar.
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando Solaar…"
if [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
  sudo add-apt-repository -y ppa:solaar-unifying/stable
  sudo apt-get update
fi
hs_install_packages solaar
hs_info "Concluído."
