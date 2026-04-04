#!/usr/bin/env bash
# OBS Studio (Flatpak) + plugin PipeWire + pacotes XWayland.
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release

obs_pipe() {
  local ver
  ver="$(curl -fsSL "https://api.github.com/repos/dimtpap/obs-pipewire-audio-capture/releases/latest" | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p')"
  [[ -n "$ver" ]] || { hs_warn "Não foi possível obter versão do plugin PipeWire; pulando."; return 0; }
  local tdir
  tdir="$(mktemp -d)"
  trap 'rm -rf "${tdir}"' RETURN
  wget -q "https://github.com/dimtpap/obs-pipewire-audio-capture/releases/download/${ver}/linux-pipewire-audio-${ver}-flatpak-30.tar.gz" -O "${tdir}/p.tar.gz" || { hs_warn "Download do plugin falhou."; return 0; }
  tar xzf "${tdir}/p.tar.gz" -C "${tdir}"
  local plugdir
  plugdir="$(hs_login_home)/.var/app/com.obsproject.Studio/config/obs-studio/plugins/linux-pipewire-audio"
  mkdir -p "${plugdir}"
  if [[ -d "${tdir}/linux-pipewire-audio" ]]; then
    cp -rf "${tdir}/linux-pipewire-audio/"* "${plugdir}/"
  fi
  sudo flatpak override --filesystem=xdg-run/pipewire-0 com.obsproject.Studio || true
}

hs_info "Instalando OBS Studio (Flatpak)…"
hs_flatpak_install_user com.obsproject.Studio
sleep 1
hs_info "Instalando dependências (PipeWire / XWayland)…"
if hs_is_arch_like || [[ "$ID" == "solus" ]]; then
  hs_install_packages wireplumber xorg-xwayland
elif hs_is_debian_like; then
  hs_install_packages wireplumber xwayland
elif hs_is_fedora_like || hs_is_suse || hs_has_ostree; then
  hs_install_packages wireplumber xorg-x11-server-Xwayland
else
  hs_install_packages wireplumber || true
fi
sleep 1
obs_pipe
disp="${DISPLAY:-:0}"
flatpak override --user \
  --socket=x11 \
  --nosocket=wayland \
  --filesystem=/tmp/.X11-unix \
  --env=QT_QPA_PLATFORM=xcb \
  --env="DISPLAY=${disp}" \
  com.obsproject.Studio || true
hs_info "OBS configurado."
