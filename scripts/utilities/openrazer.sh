#!/usr/bin/env bash
# OpenRazer + Polychromatic (Flatpak sistema).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_os_release
hs_info "Instalando OpenRazer…"
if [[ "$ID" == "ubuntu" ]] || [[ "${ID_LIKE:-}" == *ubuntu* ]]; then
  sudo add-apt-repository -y ppa:openrazer/stable
  sudo apt-get update
  hs_install_packages openrazer-meta
elif hs_is_fedora_like; then
  if hs_has_ostree; then
    if ! grep -q "plugdev" /etc/group 2>/dev/null; then
      sudo bash -c 'grep "plugdev" /lib/group >> /etc/group' || true
    fi
    cd "${TMPDIR:-/tmp}"
    fv="$(rpm -E %fedora)"
    wget "https://copr.fedorainfracloud.org/coprs/ublue-os/akmods/repo/fedora-${fv}/ublue-os-akmods-fedora-${fv}.repo"
    sudo install -o 0 -g 0 "ublue-os-akmods-fedora-${fv}.repo" "/etc/yum.repos.d/"
    rm -f "ublue-os-akmods-fedora-${fv}.repo"
    wget https://openrazer.github.io/hardware:razer.repo
    sudo install -o 0 -g 0 -m644 hardware:razer.repo /etc/yum.repos.d/hardware:razer.repo
    rm -f hardware:razer.repo
    sudo rpm-ostree refresh-md || true
    sudo rpm-ostree install -yA kmod-openrazer openrazer-daemon
  else
    hs_install_packages kernel-devel
    sudo dnf config-manager addrepo --from-repofile=https://openrazer.github.io/hardware:razer.repo
    hs_install_packages openrazer-meta
  fi
elif hs_is_suse; then
  if grep -qi "slowroll" /etc/os-release; then
    sudo zypper addrepo "https://download.opensuse.org/repositories/hardware:razer/openSUSE_Slowroll/hardware:razer.repo"
  elif grep -qi "tumbleweed" /etc/os-release; then
    sudo zypper addrepo "https://download.opensuse.org/repositories/hardware:razer/openSUSE_Tumbleweed/hardware:razer.repo"
  fi
  sudo zypper refresh
  hs_install_packages openrazer-meta
elif hs_is_arch_like; then
  hs_chaotic_aur_enable
  hs_install_packages openrazer-meta
elif [[ "$ID" == "solus" ]]; then
  hs_install_packages openrazer openrazer-current
else
  hs_die "Distribuição não mapeada (${ID})."
fi
u="$(hs_login_user)"
if [[ -n "$u" && "$u" != "root" ]]; then
  if hs_is_arch_like; then
    sudo gpasswd -a "$u" openrazer || true
  else
    sudo gpasswd -a "$u" plugdev || true
  fi
fi
hs_flatpak_install_system app.polychromatic.controller
hs_info "OpenRazer instalado. Pode ser necessário reiniciar."
