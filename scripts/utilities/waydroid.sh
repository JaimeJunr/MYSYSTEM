#!/usr/bin/env bash
# Waydroid (Wayland). Passos opcionais de libhoudini/libndk não são executados aqui — veja documentação oficial.
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
[[ "${XDG_SESSION_TYPE:-}" == "wayland" ]] || hs_die "Waydroid requer sessão Wayland (XDG_SESSION_TYPE=wayland)."
hs_info "Instalando Waydroid…"
pkgs=(waydroid python3)
if hs_is_debian_like; then
  sudo apt-get install -y curl ca-certificates
  curl -s https://repo.waydro.id | sudo bash
  pkgs+=(python3-venv)
fi
hs_install_packages "${pkgs[@]}"
sudo systemctl enable --now waydroid-container
hs_warn "Executando waydroid init (GAPPS) — pode demorar."
waydroid init -c https://ota.waydro.id/system -v https://ota.waydro.id/vendor -s GAPPS || true
if command -v firewall-cmd &>/dev/null; then
  sudo firewall-cmd --zone=trusted --add-interface=waydroid0 --permanent || true
  sudo firewall-cmd --reload || true
elif command -v ufw &>/dev/null; then
  sudo ufw allow 53 || true
  sudo ufw allow 67 || true
fi
sudo iptables -P FORWARD ACCEPT || true
hs_info "Waydroid instalado. Opcional: libhoudini/libndk via waydroid_script (veja https://docs.waydro.id)."
