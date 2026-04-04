#!/usr/bin/env bash
# Teste de cartões SD/USB (fight flash fraud).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando f3…"
hs_install_packages f3
hs_info "Concluído. Uso: f3write / f3read no dispositivo."
