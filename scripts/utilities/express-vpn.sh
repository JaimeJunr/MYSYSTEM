#!/usr/bin/env bash
# ExpressVPN — instalador universal (versão fixa no URL; verifique site para atualização).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_info "Baixando instalador ExpressVPN…"
curl -fsSLo /tmp/express-installer.run "https://www.expressvpn.works/clients/linux/expressvpn-linux-universal-4.1.1.10039.run"
chmod +x /tmp/express-installer.run
bash /tmp/express-installer.run
rm -f /tmp/express-installer.run
hs_info "Concluído."
