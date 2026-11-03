#!/usr/bin/env bash
# find_free_ips.sh - Find IPs not responding to ping/TCP in a given range.
#
# Usage: find_free_ips.sh <start_ip> <end_ip> <count>
#
# Tunables (env vars):
#   PING_TIMEOUT   - seconds to wait for an ICMP reply (default: 1)
#   PING_COUNT     - number of ICMP echo requests to send (default: 1)
#   TCP_PORTS      - space-separated ports to also probe (default: "3389")
#   TCP_TIMEOUT    - seconds to wait for a TCP connect (default: 1)
#   MAX_PARALLEL   - number of IPs probed concurrently (default: 32)

set -euo pipefail

START_IP="${1:?Usage: $0 <start_ip> <end_ip> <count>}"
END_IP="${2:?Usage: $0 <start_ip> <end_ip> <count>}"
COUNT="${3:?Usage: $0 <start_ip> <end_ip> <count>}"

if ! command -v ping &>/dev/null; then
  echo "ERROR: ping is required but not installed" >&2
  exit 1
fi

prefix="${START_IP%.*}"
start_octet="${START_IP##*.}"
end_octet="${END_IP##*.}"

PING_TIMEOUT="${PING_TIMEOUT:-1}"
PING_COUNT="${PING_COUNT:-1}"
TCP_PORTS="${TCP_PORTS-3389}"
TCP_TIMEOUT="${TCP_TIMEOUT:-1}"
MAX_PARALLEL="${MAX_PARALLEL:-32}"

# Linux ping uses "-W seconds"; macOS/BSD ping uses "-W ms". Use uname
# rather than parsing `ping -h`, which varies across platforms.
ping_timeout_arg="${PING_TIMEOUT}"
if [[ "$(uname -s)" != "Linux" ]]; then
  ping_timeout_arg=$((PING_TIMEOUT * 1000))
fi

# Resolve a timeout runner: real `timeout`, then `gtimeout` (macOS/coreutils
# via brew), then fall back to a pure-bash watcher with no external dep.
TIMEOUT_BIN=""
if command -v timeout &>/dev/null; then
  TIMEOUT_BIN="timeout"
elif command -v gtimeout &>/dev/null; then
  TIMEOUT_BIN="gtimeout"
fi

# Attempt a TCP connect with a timeout, without requiring an external
# `timeout` binary if one isn't available.
tcp_connect_with_timeout() {
  local ip="$1" port="$2" secs="$3"
  if [[ -n "${TIMEOUT_BIN}" ]]; then
    "${TIMEOUT_BIN}" "${secs}" bash -c "exec 3<>'/dev/tcp/${ip}/${port}'" 2>/dev/null
    return $?
  fi

  # Pure-bash fallback: race the connect attempt against a sleep/kill watcher.
  ( exec 3<>"/dev/tcp/${ip}/${port}" ) 2>/dev/null &
  local pid=$!
  ( sleep "${secs}"; kill -9 "${pid}" 2>/dev/null ) &
  local watcher=$!
  local rc=0
  if wait "${pid}" 2>/dev/null; then
    rc=0
  else
    rc=1
  fi
  kill -9 "${watcher}" 2>/dev/null || true
  wait "${watcher}" 2>/dev/null || true
  return "${rc}"
}

tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

probe() {
  local ip="$1"
  if ping -c "${PING_COUNT}" -W "${ping_timeout_arg}" "${ip}" &>/dev/null; then
    echo up > "${tmpdir}/${ip}"
    return
  fi
  if [[ -n "${TCP_PORTS}" ]]; then
    for port in ${TCP_PORTS}; do
      if tcp_connect_with_timeout "${ip}" "${port}" "${TCP_TIMEOUT}"; then
        echo up > "${tmpdir}/${ip}"
        return
      fi
    done
  fi
  echo down > "${tmpdir}/${ip}"
}

running=0
for ((octet = start_octet; octet <= end_octet; octet++)); do
  ip="${prefix}.${octet}"
  probe "${ip}" &
  running=$((running + 1))
  if (( running >= MAX_PARALLEL )); then
    wait -n
    running=$((running - 1))
  fi
done
wait

found=0
for ((octet = start_octet; octet <= end_octet; octet++)); do
  ip="${prefix}.${octet}"
  status="$(cat "${tmpdir}/${ip}" 2>/dev/null || echo down)"
  if [[ "${status}" == "down" ]]; then
    echo "${ip}"
    found=$((found + 1))
    if [ "${found}" -ge "${COUNT}" ]; then
      break
    fi
  fi
done