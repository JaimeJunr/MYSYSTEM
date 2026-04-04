#!/usr/bin/env bash
# Windscribe (pacote baixado do site oficial).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando Windscribe…"
cd "${TMPDIR:-/tmp}"
if hs_is_debian_like; then
  wget -O windscribe.deb "https://windscribe.com/install/desktop/linux_deb_x64"
  sudo dpkg -i windscribe.deb || sudo apt-get install -f -y
  rm -f windscribe.deb
elif hs_is_fedora_like; then
  wget -O windscribe.rpm "https://windscribe.com/install/desktop/linux_rpm_x64"
  if hs_has_ostree; then
    sudo rpm-ostree install -yA windscribe.rpm
  else
    sudo dnf install -y windscribe.rpm
  fi
  rm -f windscribe.rpm
elif hs_is_arch_like; then
  wget -O windscribe.pkg.tar.zst "https://windscribe.com/install/desktop/linux_zst_x64"
  sudo pacman -U --noconfirm windscribe.pkg.tar.zst
  rm -f windscribe.pkg.tar.zst
elif hs_is_suse; then
  wget -O windscribe.rpm "https://windscribe.com/install/desktop/linux_rpm_opensuse_x64"
  sudo zypper install -y windscribe.rpm
  rm -f windscribe.rpm
else
  hs_die "Distribuição não suportada (${ID})."
fi
hs_info "Concluído."
