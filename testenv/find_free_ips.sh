#!/usr/bin/env bash
# find_free_ips.sh - Find IPs not responding to ping in a given range.
#
# Usage: find_free_ips.sh <start_ip> <end_ip> <count>
#
# Outputs up to <count> IPs (one per line) that are down (no ping response)
# within the range [start_ip, end_ip]. Uses nmap when available (fast); falls
# back to a portable per-IP ping loop otherwise so the script does not hard-
# require nmap on the test runner.

set -euo pipefail

START_IP="${1:?Usage: $0 <start_ip> <end_ip> <count>}"
END_IP="${2:?Usage: $0 <start_ip> <end_ip> <count>}"
COUNT="${3:?Usage: $0 <start_ip> <end_ip> <count>}"

# Fast path: nmap ping-scan the range and print hosts reported as Down.
if command -v nmap &>/dev/null; then
  nmap -sn -n "${START_IP}-${END_IP##*.}" -oG - 2>/dev/null \
    | awk '/Status: Down/{print $2}' \
    | head -n "${COUNT}"
  exit 0
fi

# Fallback: ping each IP in the range and emit the ones that do not answer.
# The single-packet timeout flag differs across platforms (BSD/macOS uses -t
# seconds, Linux/iputils uses -W seconds), so pick it from uname.
case "$(uname -s)" in
  Darwin | *BSD) PING_ARGS=(-c 1 -t 1) ;;
  *)             PING_ARGS=(-c 1 -W 1) ;;
esac

prefix="${START_IP%.*}"
start_octet="${START_IP##*.}"
end_octet="${END_IP##*.}"

found=0
for ((octet = start_octet; octet <= end_octet && found < COUNT; octet++)); do
  ip="${prefix}.${octet}"
  if ! ping "${PING_ARGS[@]}" "${ip}" &>/dev/null; then
    echo "${ip}"
    found=$((found + 1))
  fi
done
