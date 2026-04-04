#!/usr/bin/env bash
# Arch-Update — bandeja e timer (Arch / CachyOS).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_is_arch_like || hs_die "Arch-Update só é suportado em distribuições Arch-like."
hs_info "Habilitando Chaotic-AUR (pacote arch-update)…"
hs_chaotic_aur_enable
hs_install_packages arch-update
systemctl --user enable --now arch-update-tray.service || true
systemctl --user enable --now arch-update.timer || true
sleep 1
arch-update --tray --enable || true
hs_info "Concluído."
