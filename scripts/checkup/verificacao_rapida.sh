#!/usr/bin/env bash
# Verificação rápida só leitura: disco, memória, carga, journal, snaps, kernels.
# Linhas ALERT:/WARN:/OK: para resumo no final.

set -uo pipefail

REAL_USER="${REAL_USER:-${USER:-}}"
REAL_HOME="${REAL_HOME:-$HOME}"

alert=0
warn=0

emit_alert() { echo "ALERT: $*"; alert=$((alert + 1)); }
emit_warn()  { echo "WARN: $*";  warn=$((warn + 1)); }
emit_ok()    { echo "OK: $*"; }

echo "=== Homestead — verificação rápida (somente leitura) ==="
echo "Utilizador: $REAL_USER"
echo ""

# --- Disco raiz ---
if read -r _ _ _ _ pct _ <<<"$(df -P / 2>/dev/null | tail -1)"; then
  pct="${pct%\%}"
  if [[ "$pct" =~ ^[0-9]+$ ]] && (( pct >= 95 )); then
    emit_alert "Disco / muito cheio (${pct}% usado)"
  elif [[ "$pct" =~ ^[0-9]+$ ]] && (( pct >= 85 )); then
    emit_warn "Disco / acima de 85% (${pct}% usado)"
  else
    emit_ok "Disco / uso ~${pct}%"
  fi
else
  emit_warn "Não foi possível ler df /"
fi

# --- Memória ---
if [[ -r /proc/meminfo ]]; then
  mapfile -t _mi < <(grep -E '^(MemTotal|MemAvailable):' /proc/meminfo | awk '{print $2}')
  if ((${#_mi[@]} >= 2)); then
    total="${_mi[0]}"
    avail="${_mi[1]}"
    if (( total > 0 )); then
      pct_free=$((100 * avail / total))
      if (( pct_free < 5 )); then
        emit_alert "RAM baixa (~${pct_free}% livre segundo MemAvailable)"
      elif (( pct_free < 12 )); then
        emit_warn "RAM apertada (~${pct_free}% livre)"
      else
        emit_ok "RAM ~${pct_free}% livre (aprox.)"
      fi
    fi
  fi
fi

# --- Carga vs CPUs ---
if [[ -r /proc/loadavg ]] && [[ -r /proc/cpuinfo ]]; then
  read -r l1 _ _ _ _ < /proc/loadavg
  cpus=$(grep -c '^processor' /proc/cpuinfo 2>/dev/null || echo 1)
  if [[ "$l1" =~ ^[0-9.]+$ ]] && [[ "$cpus" =~ ^[0-9]+$ ]] && (( cpus > 0 )); then
    st=$(awk -v l="$l1" -v c="$cpus" 'BEGIN {
      if (l > c * 2) { print "alert"; exit }
      if (l > c) { print "warn"; exit }
      print "ok"
    }')
    case "$st" in
      alert) emit_alert "Carga elevada (load1=${l1}, CPUs≈${cpus})" ;;
      warn)  emit_warn "Carga alta (load1=${l1}, CPUs≈${cpus})" ;;
      *)     emit_ok "CPU load1=${l1}, CPUs≈${cpus}" ;;
    esac
  fi
fi

# --- systemd --user failed ---
if command -v systemctl &>/dev/null; then
  failed=$(systemctl --user list-units --state=failed --no-legend 2>/dev/null | wc -l)
  failed=$(echo "$failed" | tr -d ' ')
  if [[ "$failed" =~ ^[0-9]+$ ]] && (( failed > 0 )); then
    emit_warn "systemd --user tem $failed unidade(s) em failed (systemctl --user --failed)"
  else
    emit_ok "systemd user sem unidades failed"
  fi
fi

# --- Journal em disco ---
if command -v journalctl &>/dev/null; then
  if ju=$(journalctl --disk-usage 2>/dev/null); then
    size_line=$(echo "$ju" | head -1)
    if echo "$ju" | grep -oE '[0-9]+[.,][0-9]+[MG]' | head -1 | grep -qE 'G'; then
      emit_warn "Journal grande: $size_line"
    else
      emit_ok "journal: $size_line"
    fi
  fi
fi

# --- Snaps desativados ---
if command -v snap &>/dev/null; then
  dis=$(snap list --all 2>/dev/null | grep -c disabled || true)
  dis=$(echo "$dis" | tr -d ' ')
  if [[ "$dis" =~ ^[0-9]+$ ]] && (( dis > 0 )); then
    emit_warn "$dis revisão(ões) snap desativada(s) — pode limpar com script dedicado"
  else
    emit_ok "snap sem revisões disabled listadas"
  fi
fi

# --- Flatpak (contagem, sem remover) ---
if command -v flatpak &>/dev/null; then
  unr=$(flatpak list --app 2>/dev/null | wc -l) || unr=0
  unr=$(echo "$unr" | tr -d ' ')
  emit_ok "flatpak ~$unr apps (ver unused com script dedicado)"
fi

# --- Kernels antigos (Debian/Ubuntu) ---
if command -v dpkg-query &>/dev/null; then
  imgs=$(dpkg-query -W -f='${Package}\n' 2>/dev/null | grep -c '^linux-image-[0-9]' || true)
  imgs=$(echo "$imgs" | tr -d ' ')
  running=$(uname -r 2>/dev/null || echo "?")
  if [[ "$imgs" =~ ^[0-9]+$ ]] && (( imgs > 3 )); then
    emit_warn "Muitos pacotes linux-image ($imgs); kernel em execução: $running"
  elif [[ "$imgs" =~ ^[0-9]+$ ]] && (( imgs > 0 )); then
    emit_ok "kernels: $imgs pacote(s) linux-image; em execução: $running"
  fi
fi

# --- Thumbnails ---
if [[ -d "$REAL_HOME/.cache/thumbnails" ]]; then
  sz=$(du -sh "$REAL_HOME/.cache/thumbnails" 2>/dev/null | cut -f1 || echo "?")
  emit_ok "thumbnails ~/.cache/thumbnails ~$sz"
fi

echo ""
echo "---------- Resumo ----------"
echo "ALERTAS: $alert  AVISOS: $warn"
if (( alert == 0 && warn == 0 )); then
  echo "Nenhum alerta grave detetado nesta passagem."
fi
echo "============================"
