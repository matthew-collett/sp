#!/bin/sh
set -e

BINARIES="sp sp-mcp"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

die() {
  echo "error: $1" >&2
  exit 1
}

[ -w "$INSTALL_DIR" ] || die "$INSTALL_DIR is not writable, try: sudo sh"

removed=0
for bin in $BINARIES; do
  TARGET="${INSTALL_DIR}/${bin}"
  if [ -f "$TARGET" ] && [ "$(command -v "$bin" 2>/dev/null)" = "$TARGET" ]; then
    rm -f "$TARGET"
    echo "removed ${TARGET}"
    removed=$((removed + 1))
  fi
done

[ "$removed" -gt 0 ] || echo "nothing to uninstall"
echo "done."
