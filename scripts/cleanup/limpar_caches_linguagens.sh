#!/usr/bin/env bash
# Caches de toolchains: Go, Rust/cargo, npm, pnpm, Yarn, pip, Poetry, Gradle, Maven.

REAL_USER=${REAL_USER:-${SUDO_USER:-$USER}}
REAL_HOME=${REAL_HOME:-$(getent passwd "$REAL_USER" | cut -d: -f6)}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/../lib/homestead_cleanup.sh"

echo "--- Caches de linguagens / build ---"

# Go
if command -v go &>/dev/null; then
  confirm_action "Go build cache (go clean -cache -modcache -testcache)" \
    "sudo -u \"$REAL_USER\" bash -c 'export HOME=\"$REAL_HOME\"; go clean -cache -modcache -testcache 2>/dev/null || true'" \
    "du -sh \"$REAL_HOME/.cache/go-build\" 2>/dev/null | cut -f1"
fi

# Cargo
if [ -d "$REAL_HOME/.cargo/registry" ]; then
  confirm_action "Cargo registry cache" \
    "rm -rf \"$REAL_HOME/.cargo/registry\"/*" \
    "du -sh \"$REAL_HOME/.cargo/registry\" 2>/dev/null | cut -f1"
fi

# npm
if [ -d "$REAL_HOME/.npm" ]; then
  confirm_action "npm cache" \
    "sudo -u \"$REAL_USER\" npm cache clean --force 2>/dev/null || rm -rf \"$REAL_HOME/.npm\"/*" \
    "du -sh \"$REAL_HOME/.npm\" 2>/dev/null | cut -f1"
fi

# pnpm
if [ -d "$REAL_HOME/.local/share/pnpm/store" ] || [ -d "$REAL_HOME/.pnpm-store" ]; then
  confirm_action "pnpm store prune" \
    "sudo -u \"$REAL_USER\" pnpm store prune 2>/dev/null || true" \
    "du -sch \"$REAL_HOME/.local/share/pnpm/store\" \"$REAL_HOME/.pnpm-store\" 2>/dev/null | tail -1 | cut -f1"
fi

# Yarn
if [ -d "$REAL_HOME/.cache/yarn" ]; then
  confirm_action "Yarn cache" \
    "sudo -u \"$REAL_USER\" yarn cache clean 2>/dev/null || rm -rf \"$REAL_HOME/.cache/yarn\"/*" \
    "du -sh \"$REAL_HOME/.cache/yarn\" 2>/dev/null | cut -f1"
fi

# pip
if [ -d "$REAL_HOME/.cache/pip" ]; then
  confirm_action "pip cache" \
    "sudo -u \"$REAL_USER\" pip cache purge 2>/dev/null || rm -rf \"$REAL_HOME/.cache/pip\"/*" \
    "du -sh \"$REAL_HOME/.cache/pip\" 2>/dev/null | cut -f1"
fi

# Poetry
if [ -d "$REAL_HOME/.cache/pypoetry" ]; then
  confirm_action "Poetry cache" \
    "sudo -u \"$REAL_USER\" poetry cache clear pypi --all 2>/dev/null || rm -rf \"$REAL_HOME/.cache/pypoetry\"/*" \
    "du -sh \"$REAL_HOME/.cache/pypoetry\" 2>/dev/null | cut -f1"
fi

# Gradle
if [ -d "$REAL_HOME/.gradle/caches" ]; then
  confirm_action "Gradle caches (re-download pesado)" \
    "rm -rf \"$REAL_HOME/.gradle/caches\"/*" \
    "du -sh \"$REAL_HOME/.gradle/caches\" 2>/dev/null | cut -f1"
fi

# Maven
if [ -d "$REAL_HOME/.m2/repository" ]; then
  confirm_action "Maven ~/.m2/repository (re-download pesado)" \
    "rm -rf \"$REAL_HOME/.m2/repository\"/*" \
    "du -sh \"$REAL_HOME/.m2/repository\" 2>/dev/null | cut -f1"
fi

echo "--- Concluído ---"
