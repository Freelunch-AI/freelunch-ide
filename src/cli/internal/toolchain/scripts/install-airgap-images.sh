#!/usr/bin/env bash
#
# Downloads the k3s airgap image bundle into ~/.freelunch/images.
#
# Roadmap 1.2 requires the Demo environment to come up "without touching the
# internet". A k3d cluster does not satisfy that on its own: the k3s node image
# ships *no* preloaded workload images — verified by inspecting
# /var/lib/rancher/k3s/agent/images inside both the running node and the
# pristine rancher/k3s image, which does not exist in either. Every component
# k3s installs (traefik, coredns, local-path-provisioner, klipper-lb,
# klipper-helm, metrics-server, pause, busybox) is pulled from the network as
# the cluster comes up.
#
# k3s imports any tarball it finds in /var/lib/rancher/k3s/agent/images at
# startup, so the fix is to cache the bundle on the host once and mount that
# directory into the nodes. ClusterService does the mounting; this script does
# the caching.
#
# Why the *official* bundle rather than a hand-built list: enumerating the
# running cluster's images with crictl returns six, and the release's own
# k3s-images.txt lists eight. busybox and metrics-server are pulled later than
# a fresh cluster inspection catches, so a hand-rolled list looks complete and
# fails only once the network is actually gone.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "${SCRIPT_DIR}/lib.sh"
# shellcheck source=versions.env
. "${SCRIPT_DIR}/versions.env"

fl_detect_platform
IMAGES_DIR="$(fl_images_dir)"

# These are container images: always linux, whatever the host is. The arch is
# the Docker platform's, which on Apple Silicon is arm64.
IMG_OS="linux"
bundle="k3s-airgap-images-${FL_ARCH}.tar.gz"
target="${IMAGES_DIR}/${bundle}"
expected="$(fl_expected_sum k3s-airgap-images "${K3S_VERSION}" "${IMG_OS}" "${FL_ARCH}")"

# --- already cached? --------------------------------------------------------
# The skip check is the verification: a bundle already on disk is re-hashed
# rather than trusted, so a truncated or tampered cache is replaced instead of
# being mounted into every node.
if [ -f "${target}" ] && [ "$(fl_sha256 "${target}")" = "${expected}" ]; then
  echo "k3s airgap images ${K3S_VERSION} (${FL_ARCH}) already cached at ${target}"
  exit 0
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "error: airgap image download needs curl, which is not available." >&2
  exit 1
fi

# --- download ---------------------------------------------------------------
# The release tag contains '+', which must be percent-encoded in a URL.
release_tag="${K3S_VERSION/+/%2B}"
url="https://github.com/k3s-io/k3s/releases/download/${release_tag}/${bundle}"

tmp="$(mktemp -d -t freelunch-airgap-XXXXXX)"
trap 'rm -rf "${tmp}"' EXIT

echo "downloading k3s airgap images ${K3S_VERSION} (${IMG_OS}/${FL_ARCH}); this is ~220MB"
curl --proto '=https' --tlsv1.2 --fail --show-error -sL "${url}" -o "${tmp}/${bundle}"

# --- verify before installing -----------------------------------------------
# Against checksums.txt in-tree, not a digest fetched alongside the download —
# see the note in pin-tools.sh for why that distinction matters.
actual="$(fl_sha256 "${tmp}/${bundle}")"
if [ "${actual}" != "${expected}" ]; then
  echo "error: checksum mismatch for ${bundle} at ${K3S_VERSION}" >&2
  echo "  expected ${expected}" >&2
  echo "  actual   ${actual}" >&2
  exit 1
fi

# --- install ----------------------------------------------------------------
mkdir -p "${IMAGES_DIR}"
mv "${tmp}/${bundle}" "${target}"

echo "k3s airgap images ${K3S_VERSION} cached at ${target} (sha256 verified)"
