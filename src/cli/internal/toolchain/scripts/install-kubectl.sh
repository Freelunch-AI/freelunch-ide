#!/usr/bin/env bash
#
# Installs the pinned kubectl into ~/.freelunch/bin.
#
# Why not pixi, when everything else in this repo comes from pixi:
#
#   1. A customer running `freelunch install` has Docker and nothing else. The
#      CLI has to provision kubectl itself, so pixi could only ever solve this
#      for us and not for them. Using one mechanism for both keeps the path we
#      test identical to the path we ship.
#   2. kubectl is packaged on conda-forge as `kubernetes-client`, whose newest
#      build common to every platform this workspace declares is 1.24.3 —
#      linux-aarch64 lags badly. That is years older than the k3s server we
#      run and well outside the supported version skew.
#
# Unlike k3d, we control the download here, so the digest is checked before the
# binary is ever placed on disk.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "${SCRIPT_DIR}/lib.sh"
# shellcheck source=versions.env
. "${SCRIPT_DIR}/versions.env"

fl_detect_platform
INSTALL_DIR="$(fl_install_dir)"
kubectl_bin="${INSTALL_DIR}/kubectl"
expected="$(fl_expected_sum kubectl "${KUBECTL_VERSION}")"

# --- already installed? -----------------------------------------------------
# The skip check *is* the verification — see the note in install-k3d.sh.
if [ -x "${kubectl_bin}" ] && [ "$(fl_sha256 "${kubectl_bin}")" = "${expected}" ]; then
  echo "kubectl ${KUBECTL_VERSION} already installed and verified at ${kubectl_bin}"
  exit 0
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "error: kubectl installation needs curl, which is not available." >&2
  exit 1
fi

# --- download ---------------------------------------------------------------
base="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${FL_OS}/${FL_ARCH}"
tmp="$(mktemp -d -t freelunch-kubectl-XXXXXX)"
trap 'rm -rf "${tmp}"' EXIT

echo "installing kubectl ${KUBECTL_VERSION} (${FL_OS}/${FL_ARCH}) into ${INSTALL_DIR}"
curl --proto '=https' --tlsv1.2 --fail --show-error -sL "${base}/kubectl" -o "${tmp}/kubectl"

# --- verify before installing -----------------------------------------------
# Checked against checksums.txt, not against the .sha256 published next to the
# binary. Fetching both from the same place at the same time would only prove
# the transport was intact; the pinned digest was recorded when we reviewed the
# release, so a compromised upstream republishing a matching checksum still
# fails here.
actual="$(fl_sha256 "${tmp}/kubectl")"
if [ "${actual}" != "${expected}" ]; then
  echo "error: checksum mismatch for kubectl ${KUBECTL_VERSION} (${FL_OS}/${FL_ARCH})" >&2
  echo "  expected ${expected}" >&2
  echo "  actual   ${actual}" >&2
  exit 1
fi

# --- install ----------------------------------------------------------------
mkdir -p "${INSTALL_DIR}"
chmod +x "${tmp}/kubectl"
mv "${tmp}/kubectl" "${kubectl_bin}"

echo "kubectl ${KUBECTL_VERSION} installed at ${kubectl_bin} (sha256 verified)"
