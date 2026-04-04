#!/usr/bin/env bash
# Hiddify AppImage em ~/.local (sem sudo).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"

app_dir="$(hs_login_home)/.local/bin"
desktop_dir="$(hs_login_home)/.local/share/applications"
mkdir -p "$app_dir" "$desktop_dir"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "${tmp_dir}"' EXIT
archive_file="${tmp_dir}/hiddify.tar.gz"

arch="$(uname -m)"
case "$arch" in
  x86_64) hiddify_arch="x64" ;;
  aarch64|arm64) hiddify_arch="arm64" ;;
  *) hiddify_arch="" ;;
esac

if [[ -n "$hiddify_arch" ]]; then
  latest_url="$(curl -fsSL "https://api.github.com/repos/hiddify/hiddify-app/releases/latest" | grep -Eo "https://[^\"]+Hiddify-Linux-${hiddify_arch}-AppImage\\.tar\\.gz" | head -n1)"
else
  latest_url="$(curl -fsSL "https://api.github.com/repos/hiddify/hiddify-app/releases/latest" | grep -Eo 'https://[^"]+Hiddify-Linux-[^"]+-AppImage\.tar\.gz' | head -n1)"
fi
[[ -n "$latest_url" ]] || hs_die "Não foi possível obter URL do Hiddify."

curl -fL "$latest_url" -o "$archive_file"
tar -xzf "$archive_file" -C "$tmp_dir"
extracted_appimage="$(find "$tmp_dir" -maxdepth 2 -type f -name '*.AppImage' | head -n1)"
[[ -n "$extracted_appimage" ]] || hs_die "AppImage não encontrado no arquivo."

app_file="${app_dir}/hiddify.AppImage"
install -Dm755 "$extracted_appimage" "$app_file"
desktop_file="${desktop_dir}/hiddify.desktop"
cat > "$desktop_file" << DESKTOP
[Desktop Entry]
Type=Application
Name=Hiddify
Comment=Cliente Hiddify
Exec=${app_file}
Icon=network-vpn
Terminal=false
Categories=Network;Security;
StartupNotify=true
DESKTOP
hs_info "Hiddify instalado em ${app_file}"
