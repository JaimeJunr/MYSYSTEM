#!/usr/bin/env bash
# Template Flatpak --system (nome do arquivo = ref Flathub).
set -euo pipefail
_hs_find_root() {
  local d="$1"
  while [[ "$d" != "/" ]]; do
    [[ -f "$d/go.mod" ]] && { echo "$d"; return 0; }
    d="$(dirname "$d")"
  done
  return 1
}
_here="$(cd "$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")" && pwd)"
if [[ -z "${HOMESTEAD_ROOT:-}" ]]; then
  HOMESTEAD_ROOT="$(_hs_find_root "${_here}")" || { echo "homestead: raiz do projeto não encontrada (go.mod)." >&2; exit 1; }
  export HOMESTEAD_ROOT
fi
# shellcheck source=../../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
REF="$(basename "$(readlink -f "${BASH_SOURCE[0]}")" .sh)"
[[ -n "$REF" ]] || hs_die "Referência Flatpak vazia."
hs_os_release
hs_info "Instalando via Flatpak (sistema): ${REF}"
hs_flatpak_install_system "${REF}"
hs_info "Concluído: ${REF}"
