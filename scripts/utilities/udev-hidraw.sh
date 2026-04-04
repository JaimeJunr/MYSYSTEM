#!/usr/bin/env bash
# Regras udev para dispositivo HID raw (escolha interativa no terminal).
set -euo pipefail
[[ -n "${HOMESTEAD_ROOT:-}" ]] || { echo "homestead: execute via Homestead (HOMESTEAD_ROOT)." >&2; exit 1; }
# shellcheck source=../lib/homestead_util.sh
source "${HOMESTEAD_ROOT}/scripts/lib/homestead_util.sh"
[[ -t 0 ]] || hs_die "Este script precisa de terminal interativo (use o modo sudo do Homestead, não captura de saída)."

mapfile -t usb_devices < <(lsusb)
[[ ${#usb_devices[@]} -gt 0 ]] || hs_die "Nenhum dispositivo USB encontrado."

PS3="homestead: número do dispositivo: "
select device_line in "${usb_devices[@]}"; do
  [[ -n "${device_line:-}" ]] || exit 1
  vendor_id="$(echo "$device_line" | sed -n 's/.*ID \([0-9a-f]*\):.*/\1/p')"
  [[ -n "$vendor_id" ]] || hs_die "Não foi possível obter idVendor."
  vendor_name="$(echo "$device_line" | awk '{print $7}' | tr '[:upper:]' '[:lower:]')"
  [[ -n "$vendor_name" ]] || vendor_name="usb"
  rule_num=50
  while true; do
    rule_file="/etc/udev/rules.d/${rule_num}-usb-${vendor_name}.rules"
    if [[ ! -f "$rule_file" ]]; then
      break
    fi
    rule_num=$((rule_num + 1))
  done
  echo "KERNEL==\"hidraw*\", ATTRS{idVendor}==\"${vendor_id}\", MODE=\"0666\"" | sudo tee "$rule_file" >/dev/null
  hs_info "Regra criada: ${rule_file}"
  hs_warn "Reconecte o dispositivo ou execute: sudo udevadm control --reload-rules && sudo udevadm trigger"
  break
done
