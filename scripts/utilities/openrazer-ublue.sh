#!/usr/bin/env bash
# OpenRazer em imagens Universal Blue (ujust).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
command -v ujust &>/dev/null || hs_die "Comando ujust não encontrado (use imagem Universal Blue ou o script OpenRazer padrão)."
hs_info "Executando ujust install-openrazer…"
ujust install-openrazer
hs_info "Concluído."
