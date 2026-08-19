#!/usr/bin/env bash
#
# Installs the pinned k3d into ~/.freelunch/bin.
#
# This wraps the vendored upstream installer (upstream/k3d-install.sh) rather
# than piping the published one-liner into bash. Two reasons:
#
#   1. The published command fetches install.sh from the k3d `main` branch and
#      installs whatever release is newest. Both the installer and the binary
#      are therefore unpinned, so two developers running setup on different
#      days get different clusters. Vendoring pins the installer; TAG pins the
#      binary.
#   2. `curl | bash` executes remote code we have never read. The vendored copy
#      is reviewable in-tree and changes only when we deliberately re-sync it.
#
# Re-syncing: download install.sh from the tag in versions.env, diff it against
# upstream/k3d-install.sh, and commit the result with the version bump.
#
# The same procedure must work from repository setup and from `freelunch
# install` on a customer machine, so it depends on nothing from this repository
# beyond the files sitting next to it. It lives inside the Go module so a
# future toolchain package can //go:embed the directory and ship the identical
# procedure; go:embed cannot reach outside the module. The upstream copy is
# under upstream/ rather than vendor/ because the Go tool special-cases
# directories named "vendor".

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "${SCRIPT_DIR}/lib.sh"
# shellcheck source=versions.env
. "${SCRIPT_DIR}/versions.env"

fl_detect_platform
INSTALL_DIR="$(fl_install_dir)"
k3d_bin="${INSTALL_DIR}/k3d"
expected="$(fl_expected_sum k3d "${K3D_VERSION}")"

# --- already installed? -----------------------------------------------------
# The skip check *is* the verification. Comparing the digest rather than asking
# the binary to report its own version means the only way to skip work is to
# have already proved the binary is byte-identical to what we pinned. A
# corrupted or substituted k3d that still prints the right version string does
# not pass this.
if [ -x "${k3d_bin}" ] && [ "$(fl_sha256 "${k3d_bin}")" = "${expected}" ]; then
  echo "k3d ${K3D_VERSION} already installed and verified at ${k3d_bin}"
  exit 0
fi

# --- preflight --------------------------------------------------------------
VENDORED_INSTALLER="${SCRIPT_DIR}/upstream/k3d-install.sh"
if [ ! -f "${VENDORED_INSTALLER}" ]; then
  echo "error: vendored installer missing at ${VENDORED_INSTALLER}" >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
  echo "error: k3d installation needs curl or wget, and neither is available." >&2
  exit 1
fi

# The upstream script does not create its install directory.
mkdir -p "${INSTALL_DIR}"

# Remove any existing binary before handing over to upstream.
#
# Upstream decides whether to download by running `k3d version` and comparing
# the string it prints. If we get here with a binary already present, we have
# just proved its digest is wrong — but a corrupted binary can still print the
# expected version, in which case upstream reports "already latest", downloads
# nothing, and the verification below would delete the file and fail. Deleting
# it first is what makes a bad install self-repairing instead of a hard error.
rm -f "${k3d_bin}"

# --- install ----------------------------------------------------------------
# TAG              pins the release; without it upstream resolves "latest".
# USE_SUDO=false   INSTALL_DIR is user-writable, so root is never needed.
# PATH             upstream's final testVersion step runs `command -v k3d` and
#                  exits 1 if it is not on PATH. INSTALL_DIR is deliberately
#                  not on the user's PATH, so put it there for this process
#                  only. Without this the install fails at the last step even
#                  though the binary landed correctly.
echo "installing k3d ${K3D_VERSION} (${FL_OS}/${FL_ARCH}) into ${INSTALL_DIR}"
TAG="${K3D_VERSION}" \
USE_SUDO="false" \
K3D_INSTALL_DIR="${INSTALL_DIR}" \
PATH="${INSTALL_DIR}:${PATH}" \
  bash "${VENDORED_INSTALLER}"

# --- verify -----------------------------------------------------------------
# Upstream downloads and installs in one step and checks nothing: installFile
# is documented as verifying a SHA256, but its body only chmods and copies.
# k3d does publish a checksums.txt per release, which that script simply never
# fetches. So verification happens here, after the fact, against the digest
# pinned in checksums.txt. A binary that fails is removed rather than left on
# disk for the next run to trip over.
if [ ! -x "${k3d_bin}" ]; then
  echo "error: expected k3d at ${k3d_bin} after install, but it is not there" >&2
  exit 1
fi

actual="$(fl_sha256 "${k3d_bin}")"
if [ "${actual}" != "${expected}" ]; then
  rm -f "${k3d_bin}"
  echo "error: checksum mismatch for k3d ${K3D_VERSION} (${FL_OS}/${FL_ARCH}); removed." >&2
  echo "  expected ${expected}" >&2
  echo "  actual   ${actual}" >&2
  exit 1
fi

echo "k3d ${K3D_VERSION} installed at ${k3d_bin} (sha256 verified)"
