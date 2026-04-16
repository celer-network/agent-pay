#!/bin/sh

MANUAL_ROOT_DEFAULT="${AGENTPAY_MANUAL_ROOT:-/tmp/celer_manual_test}"

basic_setup() {
    manual_root="$1"
    echo "setup testnet"
    rm -rf "${manual_root}"
    go run ${AGENTPAY}/test/manual/setup.go -logcolor -outroot "${manual_root}"
}

auto_setup() {
    manual_root="$1"
    echo "setup testnet, automatically add/approve fund and register osps"
    rm -rf "${manual_root}"
    go run ${AGENTPAY}/test/manual/setup.go -logcolor -auto -outroot "${manual_root}"
}

arg="${1}"
case ${arg} in
    auto)   auto_setup "${2:-${MANUAL_ROOT_DEFAULT}}"
            ;;
    "")     basic_setup "${MANUAL_ROOT_DEFAULT}"
            ;;
    *)      basic_setup "${arg}"
            ;;
esac
