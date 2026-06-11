#!/bin/sh
set -e

REPO="matthew-collett/sp"
BINARIES="sp sp-mcp"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

die() {
  echo "error: $1" >&2
  exit 1
}

# Detect OS
case "$(uname -s)" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *) die "unsupported OS: $(uname -s)" ;;
esac

# Detect arch
case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) die "unsupported architecture: $(uname -m)" ;;
esac

# Check dependencies
for cmd in curl tar; do
  command -v "$cmd" > /dev/null 2>&1 || die "$cmd is required but not installed"
done

# Check for a sha256 tool
if command -v sha256sum > /dev/null 2>&1; then
  sha256() { sha256sum "$1" | awk '{print $1}'; }
elif command -v shasum > /dev/null 2>&1; then
  sha256() { shasum -a 256 "$1" | awk '{print $1}'; }
else
  die "sha256sum or shasum is required but not installed"
fi

# Check install dir is writable before doing any network work
[ -w "$INSTALL_DIR" ] || die "$INSTALL_DIR is not writable, try: sudo sh"

# Resolve version
VERSION="${SP_VERSION:-}"
if [ -z "$VERSION" ]; then
  VERSION=$(curl -fsSLI -o /dev/null -w '%{url_effective}' \
    "https://github.com/${REPO}/releases/latest" | sed 's|.*/v||')
fi
[ -n "$VERSION" ] || die "could not determine latest version"

TARBALL="sp_${OS}_${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/v${VERSION}"

echo "installing sp v${VERSION} (${OS}/${ARCH})..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM HUP

# Download
curl -fsSL "${BASE_URL}/${TARBALL}" -o "${TMP_DIR}/${TARBALL}"
curl -fsSL "${BASE_URL}/checksums.txt" -o "${TMP_DIR}/checksums.txt"

# Verify checksum
EXPECTED=$(grep " ${TARBALL}$" "${TMP_DIR}/checksums.txt" | awk '{print $1}')
[ -n "$EXPECTED" ] || die "no checksum entry found for ${TARBALL}"

ACTUAL=$(sha256 "${TMP_DIR}/${TARBALL}")
[ "$EXPECTED" = "$ACTUAL" ] || die "checksum verification failed"

# Extract and install
tar -xzf "${TMP_DIR}/${TARBALL}" -C "$TMP_DIR"

for bin in $BINARIES; do
  [ -f "${TMP_DIR}/${bin}" ] || die "${bin} not found in tarball"
  cp "${TMP_DIR}/${bin}" "${INSTALL_DIR}/${bin}"
  chmod 755 "${INSTALL_DIR}/${bin}"
  echo "installed ${INSTALL_DIR}/${bin}"
done

echo "done. run 'sp --version' to verify."
