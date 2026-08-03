#!/usr/bin/env bash
# find_free_ips.sh - Find IPs not responding to ping in a given range using nmap.
#
# Usage: find_free_ips.sh <start_ip> <end_ip> <count>
#
# Outputs up to <count> IPs (one per line) that are NOT in use (no ping response)
# within the range [start_ip, end_ip]. start_ip and end_ip must share the same /24.
#
# Implementation note: `nmap -sn` only reports hosts that are *up*; it does not
# emit a "Status: Down" line for every non-responding address (especially for a
# remote subnet / when run without root). So we enumerate the full range and
# subtract the addresses nmap found to be up, rather than grepping for "Down".

set -euo pipefail

START_IP="${1:?Usage: $0 <start_ip> <end_ip> <count>}"
END_IP="${2:?Usage: $0 <start_ip> <end_ip> <count>}"
COUNT="${3:?Usage: $0 <start_ip> <end_ip> <count>}"

if ! command -v nmap &>/dev/null; then
  echo "ERROR: nmap is required but not installed" >&2
  exit 1
fi

prefix="${START_IP%.*}"        # e.g. 10.44.13
start_octet="${START_IP##*.}"  # e.g. 197
end_octet="${END_IP##*.}"      # e.g. 220

# Addresses that respond to ping (i.e. currently in use).
up_ips="$(nmap -sn -n "${START_IP}-${end_octet}" -oG - 2>/dev/null | awk '/Status: Up/{print $2}')"

found=0
for ((octet = start_octet; octet <= end_octet; octet++)); do
  ip="${prefix}.${octet}"
  if ! grep -qxF "${ip}" <<<"${up_ips}"; then
    echo "${ip}"
    found=$((found + 1))
    if [ "${found}" -ge "${COUNT}" ]; then
      break
    fi
  fi
done
