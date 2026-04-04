#!/usr/bin/env bash
# Nerd Fonts — instala JetBrains Mono e Meslo (sem interface gráfica).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
hs_info "Baixando Nerd Fonts (JetBrains Mono)…"
fonts_dir="$(hs_login_home)/.local/share/fonts/nerd-fonts"
mkdir -p "${fonts_dir}"
vers="$(curl -fsSL -o /dev/null -w '%{url_effective}' https://github.com/ryanoasis/nerd-fonts/releases/latest | awk -F/ '{print $NF}')"
[[ -n "$vers" ]] || hs_die "Não foi possível resolver versão do nerd-fonts."
curl -fsSL "https://github.com/ryanoasis/nerd-fonts/releases/download/${vers}/JetBrainsMono.tar.xz" | tar -xJf - -C "${fonts_dir}"
if command -v fc-cache &>/dev/null; then
  fc-cache -fv "$(hs_login_home)/.local/share/fonts" || true
fi
hs_info "Fontes instaladas em ${fonts_dir}"
