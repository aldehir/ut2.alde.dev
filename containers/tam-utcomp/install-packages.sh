#!/bin/bash
set -o errexit

ucc=${UT2_UCC_EXE:-ucc-bin}
system_dir=${UT2_SYSTEM_DIR:-/opt/ut2004/System}
redirect_url=${REDIRECT_URL:-https://ut2.redirect.kokuei.dev}
tmp="Downloads"

main() {
  trap cleanup EXIT

  # Install packages from the redirect
  while IFS=, read -r dest name guid; do
    download_package "$dest" "$name" "$guid"
  done <"packages.csv"
}

download_package() {
  local destination="$1"
  local name="$2"
  local guid="$3"

  mkdir -p "$tmp"
  local tmpfile="$tmp/$name.uz2"
  local decompressed="$tmp/$name"

  local download_url="$redirect_url/$name.uz2/$guid"

  echo "Downloading $download_url -> $tmpfile"
  curl -sfL -o "$tmpfile" "$download_url"
  ut2u package decompress "$tmpfile"

  echo "$decompressed -> $destination/$name"
  install -m 644 "$decompressed" "$destination/$name"

  if [[ "$name" =~ \.u$ ]]; then
    pushd "$system_dir" 1>/dev/null 2>&1
    ./$ucc exportcache $name 2>&1 | sed "s/\r/\n/"
    popd 1>/dev/null 2>&1
  fi

  echo
}

cleanup() {
  rm -rf "$tmp"
}

main
