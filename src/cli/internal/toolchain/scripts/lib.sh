#!/usr/bin/env bash
#
# Shared helpers for the toolchain installers. Sourced, never executed.

# Directory the pinned tools are installed into. User-writable by design: the
# upstream defaults are /usr/local/bin, which needs sudo and mutates the
# machine. Overridable so tests and CI can redirect it.
fl_install_dir() {
  echo "${FREELUNCH_BIN_DIR:-${HOME}/.freelunch/bin}"
}

# Sets FL_OS and FL_ARCH to the names upstream uses in its release artifacts.
fl_detect_platform() {
  FL_OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$(uname -m)" in
    x86_64 | amd64) FL_ARCH="amd64" ;;
    aarch64 | arm64) FL_ARCH="arm64" ;;
    *)
      echo "error: unsupported architecture $(uname -m)" >&2
      return 1
      ;;
  esac
  case "${FL_OS}" in
    linux | darwin) ;;
    *)
      echo "error: unsupported OS ${FL_OS}. Windows users should work inside WSL2." >&2
      return 1
      ;;
  esac
}

# sha256 of a file, as a bare hex digest. Linux ships sha256sum, macOS shasum.
fl_sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    echo "error: neither sha256sum nor shasum is available" >&2
    return 1
  fi
}

# Looks up the pinned digest for a tool/version/platform in checksums.txt.
#
# This is the anchor for both the skip check and the post-install check. It is
# deliberately read from the repository rather than fetched alongside the
# download: a checksum retrieved at install time only proves the transport was
# intact, whereas this one was recorded when the release was reviewed.
fl_expected_sum() {
  local tool="$1" version="$2" dir sum
  dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

  if [ ! -f "${dir}/checksums.txt" ]; then
    echo "error: ${dir}/checksums.txt is missing; run 'pixi run task pin:tools'" >&2
    return 1
  fi

  sum="$(awk -v t="${tool}" -v v="${version}" -v o="${FL_OS}" -v a="${FL_ARCH}" \
    '$1 == t && $2 == v && $3 == o && $4 == a {print $5}' "${dir}/checksums.txt")"

  if [ -z "${sum}" ]; then
    echo "error: no pinned checksum for ${tool} ${version} ${FL_OS}/${FL_ARCH}." >&2
    echo "       checksums.txt is stale — run 'pixi run task pin:tools' and commit it." >&2
    return 1
  fi
  echo "${sum}"
}
