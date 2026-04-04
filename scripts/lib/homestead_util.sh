#!/usr/bin/env bash
# Homestead — biblioteca compartilhada para scripts de utilidades.
# Inspirada em fluxos do projeto LinuxToys (https://linux.toys), reescrita para terminal/sudo sem Zenity.

set -euo pipefail

hs_die() {
  echo "homestead: erro: $*" >&2
  exit 1
}

hs_warn() {
  echo "homestead: aviso: $*" >&2
}

hs_info() {
  echo "homestead: $*"
}

hs_login_user() {
  printf '%s' "${REAL_USER:-${SUDO_USER:-}}"
}

hs_login_home() {
  if [[ -n "${REAL_HOME:-}" ]]; then
    printf '%s' "${REAL_HOME}"
    return
  fi
  local u
  u="$(hs_login_user)"
  if [[ -n "$u" && "$u" != "root" ]]; then
    getent passwd "$u" | cut -d: -f6
    return
  fi
  printf '%s' "${HOME}"
}

hs_os_release() {
  # shellcheck source=/dev/null
  [[ -r /etc/os-release ]] && source /etc/os-release
  : "${ID:=unknown}"
}

hs_is_debian_like() {
  [[ "${ID_LIKE:-$ID}" == *debian* || "${ID_LIKE:-$ID}" == *ubuntu* || "$ID" == "debian" || "$ID" == "ubuntu" ]]
}

hs_is_fedora_like() {
  [[ "$ID" == "fedora" || "$ID" == "rhel" || "${ID_LIKE:-}" == *fedora* || "${ID_LIKE:-}" == *rhel* ]]
}

hs_is_arch_like() {
  [[ "$ID" == "arch" || "$ID" == "cachyos" || "$ID" == "artix" || "${ID_LIKE:-}" == *arch* ]]
}

hs_is_suse() {
  [[ "$ID" =~ opensuse || "$ID" == "suse" || "${ID_LIKE:-}" == *suse* ]]
}

hs_has_ostree() {
  command -v rpm-ostree &>/dev/null
}

# Habilita repositório Chaotic-AUR (Arch e derivados). Necessário para alguns pacotes dos utilitários.
hs_chaotic_aur_enable() {
  hs_os_release
  if ! hs_is_arch_like; then
    hs_die "Chaotic-AUR só se aplica a distribuições Arch-like."
  fi
  if pacman -Slq chaotic-aur &>/dev/null; then
    return 0
  fi
  hs_info "Configurando Chaotic-AUR…"
  sudo sed -i '/\[chaotic-aur\]/,/Include = \/etc\/pacman.d\/chaotic-mirrorlist/ d' /etc/pacman.conf 2>/dev/null || true
  if ! sudo pacman-key --recv-key 3056513887B78AEB --keyserver keyserver.ubuntu.com \
    || ! sudo pacman-key --lsign-key 3056513887B78AEB; then
    hs_die "Falha ao adicionar chaves Chaotic-AUR."
  fi
  if ! sudo pacman -U --noconfirm \
    'https://cdn-mirror.chaotic.cx/chaotic-aur/chaotic-keyring.pkg.tar.zst' \
    'https://cdn-mirror.chaotic.cx/chaotic-aur/chaotic-mirrorlist.pkg.tar.zst'; then
    hs_die "Falha ao instalar keyring Chaotic-AUR."
  fi
  printf '\n[chaotic-aur]\nInclude = /etc/pacman.d/chaotic-mirrorlist\n' | sudo tee -a /etc/pacman.conf >/dev/null
  sudo pacman -Syy
  pacman -Slq chaotic-aur &>/dev/null || hs_die "Chaotic-AUR não ficou disponível após configurar."
}

# Instala RPM Fusion (Fedora) quando necessário para codecs (ex.: VLC).
hs_rpmfusion_enable() {
  hs_os_release
  local fedora_version
  fedora_version="$(rpm -E %fedora 2>/dev/null)" || hs_die "Não foi possível detectar versão do Fedora."
  [[ "$fedora_version" == "%fedora" ]] && hs_die "Versão Fedora inválida para RPM Fusion."
  local install_cmd=(sudo dnf install -y)
  if hs_has_ostree; then
    install_cmd=(sudo rpm-ostree install -yA)
  fi
  if ! rpm -qi rpmfusion-free-release &>/dev/null; then
    "${install_cmd[@]}" "https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-${fedora_version}.noarch.rpm"
  fi
  if ! rpm -qi rpmfusion-nonfree-release &>/dev/null; then
    "${install_cmd[@]}" "https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-${fedora_version}.noarch.rpm"
  fi
}

# Instala pacotes nativos conforme o gestor detectado (sem Zenity / sem paru interativo).
hs_install_packages() {
  hs_os_release
  local pkgs=("$@")
  [[ ${#pkgs[@]} -gt 0 ]] || return 0

  if hs_has_ostree; then
    local p
    for p in "${pkgs[@]}"; do
      rpm -qi "$p" &>/dev/null || sudo rpm-ostree install -yA "$p"
    done
    return 0
  fi

  case "${ID_LIKE:-$ID}" in
    *debian*|*ubuntu*)
      local pak
      for pak in "${pkgs[@]}"; do
        dpkg -s "$pak" &>/dev/null || sudo apt-get install -y "$pak"
      done
      ;;
    *arch*|*archlinux*|*cachyos*)
      local pak
      for pak in "${pkgs[@]}"; do
        pacman -Qi "$pak" &>/dev/null || sudo pacman -S --noconfirm "$pak"
      done
      ;;
    *rhel*|*fedora*|*centos*)
      local pak
      for pak in "${pkgs[@]}"; do
        rpm -qi "$pak" &>/dev/null || sudo dnf install -y "$pak"
      done
      ;;
    *suse*)
      local pak
      for pak in "${pkgs[@]}"; do
        rpm -qi "$pak" &>/dev/null || sudo zypper install -y "$pak"
      done
      ;;
    *solus*)
      local pak
      for pak in "${pkgs[@]}"; do
        eopkg list-installed | grep -qx "$pak" || sudo eopkg install -y "$pak"
      done
      ;;
    *)
      hs_die "Distribuição não suportada para instalação automática de pacotes (${ID})."
      ;;
  esac
}

# Garante flatpak + remoto flathub (usuário e sistema).
hs_flatpak_ensure() {
  hs_os_release
  if ! command -v flatpak &>/dev/null; then
    hs_info "Instalando Flatpak…"
    case "${ID_LIKE:-$ID}" in
      *debian*|*ubuntu*) sudo apt-get update && sudo apt-get install -y flatpak ;;
      *arch*|*archlinux*|*cachyos*) sudo pacman -S --noconfirm flatpak ;;
      *suse*) sudo zypper install -y flatpak ;;
      *fedora*|*rhel*) sudo dnf install -y flatpak ;;
      *solus*) sudo eopkg install -y flatpak ;;
      *) hs_die "Instale flatpak manualmente nesta distribuição." ;;
    esac
  fi
  flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
  sudo flatpak remote-add --if-not-exists --system flathub https://dl.flathub.org/repo/flathub.flatpakrepo
}

hs_flatpak_install_user() {
  local ref="$1"
  hs_flatpak_ensure
  flatpak install --or-update --user --noninteractive flathub "${ref}"
}

hs_flatpak_install_system() {
  local ref="$1"
  hs_flatpak_ensure
  sudo flatpak install --or-update --system --noninteractive flathub "${ref}"
}
