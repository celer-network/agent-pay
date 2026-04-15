#!/bin/sh
set -eu

if [ -z "${AGENTPAY:-}" ]; then
	echo "AGENTPAY must point to the agent-pay repo root" >&2
	exit 1
fi

MANUAL_ROOT="${AGENTPAY_MANUAL_ROOT:-/tmp/celer_manual_test}"
MANUAL_ROOT="${MANUAL_ROOT%/}"
LOG_DIR="${MANUAL_ROOT}-logs"
CLI="${AGENTPAY}/bin/osp-cli"
PAY_ID=""
SETUP_PID=""
OSP1_PID=""
OSP2_PID=""

wait_for_log() {
	log_file="$1"
	needle="$2"
	timeout_sec="$3"
	count=0
	while [ "$count" -lt "$timeout_sec" ]; do
		if [ -f "$log_file" ] && grep -q "$needle" "$log_file"; then
			return 0
		fi
		sleep 1
		count=$((count + 1))
	done
	return 1
}

cleanup() {
	status=$?
	set +e
	if [ -n "$OSP1_PID" ]; then
		kill "$OSP1_PID" 2>/dev/null || true
	fi
	if [ -n "$OSP2_PID" ]; then
		kill "$OSP2_PID" 2>/dev/null || true
	fi
	if [ -n "$SETUP_PID" ]; then
		kill "$SETUP_PID" 2>/dev/null || true
	fi
	pkill -f "${MANUAL_ROOT}/chaindata" 2>/dev/null || true
	pkill -f "${MANUAL_ROOT}/profile/o1_profile.json" 2>/dev/null || true
	pkill -f "${MANUAL_ROOT}/profile/o2_profile.json" 2>/dev/null || true
	if [ "$status" -ne 0 ]; then
		echo "manual smoke failed; logs kept in ${LOG_DIR}" >&2
		if [ -n "$PAY_ID" ]; then
			echo "last payment id: ${PAY_ID}" >&2
		fi
	fi
	exit "$status"
}

trap cleanup EXIT INT TERM

rm -rf "$LOG_DIR"
mkdir -p "$LOG_DIR"

cd "$AGENTPAY"

echo "building osp-cli"
go build -o "$CLI" ./tools/osp-cli

echo "starting manual setup"
AGENTPAY="$AGENTPAY" AGENTPAY_MANUAL_ROOT="$MANUAL_ROOT" ./test/manual/setup.sh auto >"${LOG_DIR}/setup.log" 2>&1 &
SETUP_PID=$!
if ! wait_for_log "${LOG_DIR}/setup.log" "Local testnet setup finished." 120; then
	echo "manual setup did not become ready" >&2
	tail -n 80 "${LOG_DIR}/setup.log" >&2 || true
	exit 1
fi

echo "starting OSP1 and OSP2"
AGENTPAY="$AGENTPAY" AGENTPAY_MANUAL_ROOT="$MANUAL_ROOT" ./test/manual/run_osp.sh 1 >"${LOG_DIR}/o1.log" 2>&1 &
OSP1_PID=$!
AGENTPAY="$AGENTPAY" AGENTPAY_MANUAL_ROOT="$MANUAL_ROOT" ./test/manual/run_osp.sh 2 >"${LOG_DIR}/o2.log" 2>&1 &
OSP2_PID=$!
if ! wait_for_log "${LOG_DIR}/o1.log" "admin HTTP: localhost:8190" 60; then
	echo "OSP1 did not become ready" >&2
	tail -n 80 "${LOG_DIR}/o1.log" >&2 || true
	exit 1
fi
if ! wait_for_log "${LOG_DIR}/o2.log" "admin HTTP: localhost:8290" 60; then
	echo "OSP2 did not become ready" >&2
	tail -n 80 "${LOG_DIR}/o2.log" >&2 || true
	exit 1
fi

echo "registering inter-OSP stream"
"$CLI" -adminhostport localhost:8190 -registerstream -peer 00290a43e5b2b151d530845b2d5a818240bc7c70 -peerhostport localhost:10002 >/tmp/agentpay_manual_registerstream.out 2>&1

echo "opening OSP-to-OSP channel"
"$CLI" -adminhostport localhost:8190 -openchannel -peer 00290a43e5b2b151d530845b2d5a818240bc7c70 -selfdeposit 10 -peerdeposit 10 >/tmp/agentpay_manual_openchannel.out 2>&1
sleep 2

echo "building routing tables"
curl -sS -X POST -H 'Content-Type: application/json' -d '{}' http://localhost:8190/admin/route/build >/tmp/agentpay_manual_route_o1.json
curl -sS -X POST -H 'Content-Type: application/json' -d '{}' http://localhost:8290/admin/route/build >/tmp/agentpay_manual_route_o2.json

echo "sending payment"
send_output=$("$CLI" -adminhostport localhost:8190 -sendtoken -receiver 00290a43e5b2b151d530845b2d5a818240bc7c70 -amount 0.01 2>&1)
printf '%s\n' "$send_output"
PAY_ID=$(printf '%s\n' "$send_output" | sed -n 's/.*requested to send payment \([0-9a-f][0-9a-f]*\).*/\1/p' | tail -n 1)
if [ -z "$PAY_ID" ]; then
	echo "failed to parse payment id from sendtoken output" >&2
	exit 1
fi

echo "verifying payment state"
pay_output=$("$CLI" -profile "${MANUAL_ROOT}/profile/o2_profile.json" -storedir "${MANUAL_ROOT}/store/00290a43e5b2b151d530845b2d5a818240bc7c70" -dbview pay -payid "$PAY_ID")
printf '%s\n' "$pay_output"
printf '%s\n' "$pay_output" | grep -q 'state COSIGNED_PAID'

echo "verifying channel balances"
channel_output=$("$CLI" -profile "${MANUAL_ROOT}/profile/o1_profile.json" -storedir "${MANUAL_ROOT}/store/0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e" -dbview channel -peer 00290a43e5b2b151d530845b2d5a818240bc7c70)
printf '%s\n' "$channel_output"
printf '%s\n' "$channel_output" | grep -q 'self free balance: 9990000000000000000'
printf '%s\n' "$channel_output" | grep -q 'peer free balance: 10010000000000000000'

echo "manual smoke passed"
echo "payment id: ${PAY_ID}"
echo "logs: ${LOG_DIR}"